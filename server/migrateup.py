import sqlite3
import bcrypt

conn = sqlite3.connect("database.db")
print("Connected to database successfully")

c = conn.cursor()

c.execute(
    """CREATE TABLE users (
    username text,
    first_name text,
    last_name text,
    email text,
    password blob,
    salary integer
    )"""
)

salt = "$2b$12$gBL4O3YeTVAbNSviFoOl2e".encode()
pass1 = "helloworld"
pass2 = "12345"
hash1 = bcrypt.hashpw(pass1.encode(), salt)
hash2 = bcrypt.hashpw(pass2.encode(), salt)
dummy_accounts = [
    ("djohnoe", "John", "Doe", "johndoe@gmail.com", hash1, 1000),
    ("janeDOE", "Jane", "Doe", "jane@gmail.com", hash2, 900000),
]
c.executemany("INSERT INTO users VALUES (?,?,?,?,?,?)", dummy_accounts)

### dev only

###

conn.commit()
conn.close()
