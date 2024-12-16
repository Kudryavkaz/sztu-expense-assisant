import os
import json
import logging
import requests
from datetime import datetime
from flask import Flask, render_template, request, redirect, url_for
from pandas import DataFrame
from typing import Tuple

BACKEND_SERVER_HOST = os.environ.get("BACKEND_SERVER_HOST", "localhost")
BACKEND_SERVER_PORT = int(os.environ.get("BACKEND_SERVER_PORT", "3000"))


def pull_info(student_id: str, student_passwd: str) -> None:
    payload = json.dumps({
        "sztu_account": student_id,
        "sztu_password": student_passwd,
    })
    headers = {
        'Content-Type': 'application/json',
        'Accept': '*/*',
        'Host': f"{BACKEND_SERVER_HOST}:{BACKEND_SERVER_PORT}",
    }
    response = requests.get(f"http://{BACKEND_SERVER_HOST}:{BACKEND_SERVER_PORT}/v1/expense/data",
                            data=payload, headers=headers)
    if response.status_code != 200:
        raise Exception(f"Failed to pull data, status code: {response.status_code}")

    return


def get_table(student_id: str, page: int, per_page: int) -> Tuple[DataFrame, int]:
    body = json.dumps({
        "sztu_account": student_id,
        "page": page,
        "per_page": per_page
    })
    headers = {
        'Content-Type': 'application/json',
        'Accept': '*/*',
        'Host': f"{BACKEND_SERVER_HOST}:{BACKEND_SERVER_PORT}",
    }
    try:
        response = requests.get(
            f"http://{BACKEND_SERVER_HOST}:{BACKEND_SERVER_PORT}/v1/expense/table", data=body, headers=headers)
        if response.status_code != 200:
            raise Exception(f"Failed to get table data, status code: {response.status_code}")

        data = json.loads(response.text)
        total = data['data']['total']
        df = DataFrame(data['data']['expense_info_list'])
    except Exception as e:
        logging.error(e)
        raise e

    return df, total


def get_daily_expenditure(student_id: str) -> DataFrame:
    body = json.dumps({
        "sztu_account": student_id,
        "start_time": int(datetime.strptime("2024-08-01", "%Y-%m-%d").timestamp() * 1000),
        "end_time": int(datetime.now().timestamp() * 1000)
    })
    headers = {
        'Content-Type': 'application/json',
        'Accept': '*/*',
        'Host': f"{BACKEND_SERVER_HOST}:{BACKEND_SERVER_PORT}",
    }
    try:
        response = requests.get(
            f"http://{BACKEND_SERVER_HOST}:{BACKEND_SERVER_PORT}/v1/expense/timeline", data=body, headers=headers)
        if response.status_code != 200:
            raise Exception(f"Failed to get table data, status code: {response.status_code}")

        data = json.loads(response.text)
        df = DataFrame(data['data']['expense_info_list'])
    except Exception as e:
        logging.error(e)
        raise e

    return df


def show(debug):
    app = Flask(__name__, template_folder='../templates', static_folder='../static')

    # 主界面路由
    @app.route('/', methods=['GET', 'POST'])
    def index():
        if request.method == 'POST':
            student_id = request.form['student_id']
            return redirect(url_for('user', student_id=student_id))
        return render_template('index.html')

    # 更新数据接口
    @app.route('/update/<student_id>', methods=['POST'])
    def update_data(student_id):
        student_passwd = request.form['new_data']

        pull_info(student_id, student_passwd)

        return redirect(url_for('user', student_id=student_id))

    # 用户界面路由
    @app.route('/user/<student_id>')
    def user(student_id):

        page = request.args.get('page', 1, type=int)
        per_page = 10

        df, total_records = get_table(student_id, page, per_page)

        if df is not None and not df.empty:
            total_pages = (total_records + per_page - 1) // per_page

            daily_expenditure = get_daily_expenditure(student_id)
            df.columns = ['用户操作', '时间', '地点', '金额']
            tables = [df.to_html(classes='table table-striped', header="true", index=False)]  # type: ignore

            daily_expenditure.columns = ['day', 'total_expenditure']

            return render_template(
                'user.html',
                tables=tables,
                student_id=student_id,
                page=page,
                total_pages=total_pages,
                daily_expenditure=daily_expenditure.to_dict(orient='records')
            )
        else:
            return render_template('user.html', student_id=student_id,
                                   message="No data found for student_id: {}".format(student_id))

    app.run(host="0.0.0.0", port=5000, debug=debug)
