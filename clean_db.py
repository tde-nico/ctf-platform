import sqlite3

db = sqlite3.connect('database.db')

db.execute("DELETE FROM submissions")
print('Cleaned submissions')
db.execute("DELETE FROM solves")
print('Cleaned solves')
db.execute("DELETE FROM users WHERE is_admin=0")
print('Cleaned NON admin users')
