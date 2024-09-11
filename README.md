# Image Labeling

 This is a set of scripts for images labeling using machine learning tools. It has a client server architecture. A small Python script that runs =microsoft/Florence2= model on a GPU cloud machine. And local bash script that fetch and populate detailed image Description.

- Get GPU instance on vast.ai RTX 3060 should do
- Run `push.sh` script to copy scripts to remote machine
- SSH and map local port 5000 `ssh -p 33526 root@174.95.30.134 -L 5000:localhost:5000`
- Run `./setup.sh` on remote machine. This will install Python dependencies.
- Run `./caption.py` on remote machine. It will download the Florence2 model and start HTTP server on `localhost:5000`. The server `/upload` endpoint listen for POST request with an image body and reply with detailed image caption in JSON format.
- Now from local machine run `lable.sh <path_to_img_dir>` this will recursively find all JPG images in that directory and start getting captures for them and storing it to XMP metadata  as `Description`.