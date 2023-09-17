from flask import Flask, request
import sqlite3
import json

app = Flask(__name__)


def get_db():
    conn = sqlite3.connect("database.db")
    return conn


@app.route("/accounts/signin")
def signin():
    return "hi"


@app.route("/accounts/signup", methods=["POST"])
def signup():
    conn = get_db()
    c = conn.cursor()

    body = json.loads(request.data)
    print(body)

    c.execute(
        "SELECT rowid, * FROM users WHERE username = ? OR email = ?",
        [body["username"], body["email"]],
    )
    users = c.fetchall()
    print(users)
    p

    conn.close()
    return {"ok": True}
