import sqlite3
import os

db = sqlite3.connect('database.db')

configs = [
	{
		'key': 'registration-allowed',
		'type': 'bool',
		'value': '1',
		'desc': 'Allow new users to register',
	}, {
		'key': 'chall-min-points',
		'type': 'int',
		'value': '100',
		'desc': 'Challenge minimum points',
	}, {
		'key': 'chall-points-decay',
		'type': 'int',
		'value': '18',
		'desc': 'Challenge solves needed to fully decay',
	}, {
		'key': 'telegram-bot-enable',
		'type': 'bool',
		'value': '1',
		'desc': 'Notify first blood on Telegram',
	}, {
		'key': 'telegram-bot-chat',
		'type': 'text',
		'value': '',
		'desc': 'Telegram chat ID',
	}
]

for i in range(len(configs)):
	db.execute("INSERT INTO config values (?, ?, ?, ?)", [
		configs[i]['key'],
		configs[i]['type'],
		configs[i]['value'],
		configs[i]['desc']
	])
	db.commit()

db.execute("INSERT INTO keys values (?, ?)", [
	'secret-key',
	os.urandom(32).hex(),
])
db.execute("INSERT INTO keys values (?, ?)", [
	'telegram--key',
	'',
])
db.commit()
