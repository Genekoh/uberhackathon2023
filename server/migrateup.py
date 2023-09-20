import sqlite3
import bcrypt

conn = sqlite3.connect("database.db")
print("Connected to database successfully")

c = conn.cursor()

c.execute(
    """CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE,
    firstname TEXT,
    lastname TEXT,
    email TEXT UNIQUE,
    passwordhash BLOB,
    salary INTEGER,
    accountlevel INTEGER
    )"""
)

c.execute(
    """CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BLOB NOT NULL,
	expiry REAL NOT NULL
)"""
)
c.execute("CREATE INDEX sessions_expiry_idx ON sessions(expiry)")

# bookings {
#   userid int
#   location smth
#   carpool refrence
#   createdAt time
#   expiresAt time
# }

# carpool {
# createdAt
# number int
# }
c.execute(
    """CREATE TABLE carpools (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    createdAt INTEGER,
    expiresAt INTEGER,
    size INTEGER
  )"""
)

c.execute(
    """CREATE TABLE bookings (
    userid INTEGER ,
    carpoolid INTEGER,
    pickuplat REAL,
    pickuplon REAL,
    destlat REAL,
    destlon REAL,
    createdAt INTEGER,
    expiresAt INTEGER,
    cost REAL,
    FOREIGN KEY(userid) REFERENCES users(id) ON DELETE CASCADE
    FOREIGN KEY(carpoolid) REFERENCES carpools(id) ON DELETE CASCADE
)"""
)

### dev only
salt = "$2b$12$gBL4O3YeTVAbNSviFoOl2e".encode()
pass1 = "helloworld"
pass2 = "12345"
hash1 = bcrypt.hashpw(pass1.encode(), salt)
hash2 = bcrypt.hashpw(pass2.encode(), salt)
dummy_accounts = [
    ("djohnoe", "John", "Doe", "johndoe@gmail.com", hash1, 1000, 0),
    ("janeDOE", "Jane", "Doe", "jane@gmail.com", hash2, 90000, 1),
]
c.executemany(
    "INSERT INTO users (username, firstname, lastname, email, passwordhash, salary, accountlevel) VALUES (?,?,?,?,?,?,?)",
    dummy_accounts,
)
###

conn.commit()
conn.close()
print("Migrated up")
