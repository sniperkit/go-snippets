import logging
import sys
import json
import os

log = logging.getLogger()
log.setLevel(10)
consoleHandler = logging.StreamHandler(sys.stdout)
log.addHandler(consoleHandler)

def runAlgorithm(x, y, outputDir):
    z = int(x) + int(y)
    log.info(str(x)+" + "+str(y)+" = " +str(z))

    if not outputDir == "":
        fname = os.path.join(outputDir, str(x) + "_" + str(y) + "_output.txt")
        if os.path.exists(fname):
            append_write = 'a'
        else:
            append_write = 'w+'
        target = open(fname, 'w+')
        target.write(str(x)+" + "+str(y)+" = " +str(z))
        target.write("\n")
        target.close()

    # write to seed manifest
    if not outputDir == "":
        target = open(os.path.join(outputDir, "seed.outputs.json"), 'a')
        target.write("\"x\": " + str(x) + ",\n")
        target.write("\"y\": " + str(y) + ",\n")
        target.write("\"total\": " + str(z))
        target.close()
    return z

if __name__ == '__main__':
    sys_stdout = sys.stdout
    argv = sys.argv
    if argv is None:
        log.error('No inputs passed to job')
        sys.exit(2)
    argc = len(argv) - 1

    # Must always have an input file
    inputStr = argv[1]

    # output file is optional
    outputStr = ""
    if len(argv) > 2:
        outputStr = argv[2]
        if not os.path.exists(outputStr):
            os.makedirs(outputStr)

    input_obj = open(inputStr, "r")
    inputs = input_obj.readlines()

    if not outputStr == "":
        fname = os.path.join(outputStr, "seed.outputs.json")
        if os.path.exists(fname):
            target = open(fname, 'w')
        else:
            target = open(fname, 'w+')

        target.write('{\n')
        target.flush()
        target.close()

    for idx, line in enumerate(inputs):
        xy = line.split()
        total = runAlgorithm(xy[0], xy[1], outputStr)
        if idx < len(inputs)-1 and outputStr != "":
            target = open(fname, 'a')
            target.write(",\n")
            target.close()
        else:
            target = open(fname, 'a')
            target.write("\n")
            target.flush()
            target.close()

    if not outputStr == "":
        target = open(os.path.join(outputStr, "seed.outputs.json"), 'a')
        target.write("}")
        target.close

    log.info('Completed Python Wrapper')

    sys.exit(0)
