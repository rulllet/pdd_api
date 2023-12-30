import os


directory = "./static/img"
files = os.listdir(directory)

for i in files:
    if i not in ['image.jpg', 'favicon.ico']:
        name = i.split("-")
        name[1] = str(int(name[1]) + 1)
        name = "-".join(name)
        print(name)
        os.rename(i, name)
