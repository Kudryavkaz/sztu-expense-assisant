import time
from math import ceil
from typing import List

from pydantic import BaseModel, Field, SecretStr
from tencentcloud.common.credential import Credential
from tencentcloud.common.exception import TencentCloudSDKException
from tencentcloud.ocr.v20181119 import models
from tencentcloud.ocr.v20181119.ocr_client import OcrClient
from src.log import logger
from src.config import TENCENTCLOUD_SECRET_ID, TENCENTCLOUD_SECRET_KEY


def _detect_text_to_lines(text_detections: List[models.TextDetection]) -> List[str]:
    lines: List[str] = []
    for idx in range(1, len(text_detections)):
        pre_dct, cur_dct = text_detections[idx - 1], text_detections[idx]

        font_width = ceil((cur_dct.ItemPolygon.Width + pre_dct.ItemPolygon.Width) / (len(cur_dct.DetectedText) + len(pre_dct.DetectedText)))  # type: ignore

        # cmp horizontal position
        if cur_dct.ItemPolygon.X < pre_dct.ItemPolygon.X:  # type: ignore
            # swap
            pre_dct, cur_dct = cur_dct, pre_dct

        if abs(cur_dct.ItemPolygon.Y - pre_dct.ItemPolygon.Y) < 5:  # type: ignore
            # same line
            num_space = (cur_dct.Polygon[0].X - pre_dct.Polygon[1].X) // font_width  # type: ignore
            text_detections[idx].DetectedText = pre_dct.DetectedText + num_space * " " + cur_dct.DetectedText  # type: ignore
            text_detections[idx - 1].DetectedText = ""

    lines = [dct.DetectedText for dct in text_detections if dct.DetectedText]  # type: ignore

    return lines


class ImageRecognizer(BaseModel):
    tencent_cloud_secret_id: SecretStr = Field(alias="secret_id", default=None)  # type: ignore
    tencent_cloud_secret_key: SecretStr = Field(alias="secret_key", default=None)  # type: ignore
    tencent_cloud_region: str = Field(default="ap-guangzhou", alias="region")

    workers: int = 5
    language: str = "zh"
    max_retries: int = 5

    @property
    def credential(self):
        if not self.tencent_cloud_secret_id:
            raise ValueError("Secret id is empty.")
        if not self.tencent_cloud_secret_key:
            raise ValueError("Secret key is empty.")

        return Credential(self.tencent_cloud_secret_id.get_secret_value(), self.tencent_cloud_secret_key.get_secret_value())

    @property
    def client(self):
        return OcrClient(self.credential, "ap-guangzhou")

    def recognize(self, base64_image: str) -> str:
        request = models.GeneralAccurateOCRRequest()
        request.ImageBase64 = base64_image

        response = None

        retry_sec = 1
        for _ in range(self.max_retries):
            try:
                response = self.client.GeneralAccurateOCR(request)
            except TencentCloudSDKException as e:
                if e.code == "FailedOperation.ImageNoText":
                    logger.warning("[TencentCloud OCR] No text detected in image.")
                    return ""
                elif e.code == "FailedOperation.ImageDecodeFailed":
                    logger.warning("[TencentCloud OCR] Image decode error.")
                    return ""
                logger.error(f"[TencentCloud OCR] Sdk error: {e}. retry after {retry_sec} seconds.")
                time.sleep(retry_sec)
                retry_sec <<= 1
            except Exception as e:
                logger.error(f"[TencentCloud OCR] Request error: {e}.")
            else:
                break

        if not response:
            raise ValueError(f"Failed to recognize image after {self.max_retries} retries.")
        if not isinstance(response.TextDetections, list):
            raise ValueError(f"Invalid response: {response}")

        sorted_text_detections = sorted(response.TextDetections, key=lambda x: x.ItemPolygon.Y)

        lines = _detect_text_to_lines(sorted_text_detections)

        return "\n".join([line for line in lines if line])


def ocr(base64_image) -> str:
    image_recognizer = ImageRecognizer(
        secret_id=TENCENTCLOUD_SECRET_ID,  # type: ignore
        secret_key=TENCENTCLOUD_SECRET_KEY,  # type: ignore
    )
    recognized_content = image_recognizer.recognize(base64_image)  # type: ignore
    return recognized_content
