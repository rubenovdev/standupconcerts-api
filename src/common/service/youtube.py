from flask import Flask
from pytube import YouTube
from flask import request

app = Flask(__name__)

@app.route("/upload-video", methods=["POST"])
def getVideo():
    request_data = request.json
    
    link = request_data["link"]
    outDir = request_data["outDir"]
    filename = request_data["filename"]

    YouTube(link).streams.filter(progressive=True, file_extension='mp4').order_by(
        'resolution').desc().first().download(output_path=outDir, filename=filename)
    return "successfully"

if __name__ == '__main__':
    # run app in debug mode on port 5000
    app.run(debug=True, port=5000)