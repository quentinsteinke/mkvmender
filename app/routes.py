from app import app
from flask import render_template, request, redirect, url_for

@app.route('/', methods=['GET', 'POST'])
def index():
    if request.method == 'POST':
        file_name = request.form.get('file_name')
        return redirect(url_for('submit', file_name=file_name))
    
    return render_template('index.html')

@app.route('/submit/<file_name>')
def submit(file_name):
    return render_template('submit.html', file_name=file_name)
