#coding="utf-8"


import os
from moviepy.editor import VideoFileClip, concatenate_videoclips
        
def batch_volumex(path, x):
    # 函数功能：在指定路径下，将该文件夹的视频声音调为x倍
    origin_path = os.getcwd()
    os.chdir(path)
    out_dir="./enhance"
    if not os.path.exists(out_dir):
        os.makedirs(out_dir)
    for fname in os.listdir():
        if not fname.endswith(".mp4"):
            continue
        clip = VideoFileClip(fname)
        newclip = clip.volumex(x)
        newclip.write_videofile(out_dir+"/"+fname)
        print(fname)
        # break
    os.chdir(origin_path)

path = 'D:/download/加密与安全'
batch_volumex(path, 10) # 音调提高10分贝