import cairosvg
import scrython
import urllib3
import requests
import time
import glob
import argparse
parser = argparse.ArgumentParser()

parser.add_argument("-o", "--output_directory", default="./set_symbols/", help="Output directory")

args = parser.parse_args()
sets = scrython.Sets()
existing_sets = glob.glob(args.output_directory+"*.png")
for s in sets.data():
    if args.output_directory+s["code"]+".png" in existing_sets: # We strip the .png suffix, which should be there because of the glob
        print(s["name"], "already has a file in", s["code"]+".png", "- we'll skip it.")
        continue

    print("Saving", s["name"], "from", s["icon_svg_uri"])
    resp = requests.get(s["icon_svg_uri"])
    resp.raise_for_status()
    svg = cairosvg.svg2png(bytestring=resp.text.encode("utf-8"), write_to=args.output_directory+s["code"]+".png")
    time.sleep(0.1) # Be kind to the endpoint
