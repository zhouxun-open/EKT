#!/usr/bin/env python
# -*- coding:utf-8 -*-

import os
import glob

def pwd():
    return os.path.dirname(os.path.abspath(__file__))

def get_containers():
    cs = []
    for conf in glob.glob('*.genesis.json'):
        cs.append(conf.strip('.genesis.json'))
    return cs

def pull(image):
    os.system('docker pull %s' % image)

def run(image):
    cs = get_containers()
    for name in cs:
        conf = os.path.join(pwd(), '%s.genesis.json' % name)
        netconf = os.path.join(pwd(), 'localnet.json')
        port = name.split('_')[1]
        cmd = 'docker run -d --restart always --name %s -v %s:/root/genesis.json -v %s:/root/localnet.json -p %s:%s %s' % (name, conf, netconf, port, port, image)
        os.system(cmd)

def start():
    cs = get_containers()
    for name in cs:
        os.system('docker start %s' % name)

def stop():
    cs = get_containers()
    for name in cs:
        os.system('docker stop %s' % name)

def clean():
    cs = get_containers()
    for name in cs:
        os.system('docker rm -f %s' % name)

def restart():
    cs = get_containers()
    for name in cs:
        os.system('docker restart %s' % name)

def date():
    cs = get_containers()
    for name in cs:
        os.system('docker exec %s date' % name)

def exec_cmd(cmd):
    cs = get_containers()
    for name in cs:
        os.system('docker exec %s %s' % (name, cmd))


if __name__ == '__main__':
    import sys
    import inspect
    if len(sys.argv) < 2:
        print "Usage:"
        for k, v in sorted(globals().items(), key=lambda item: item[0]):
            if inspect.isfunction(v) and k[0] != "_":
                args, __, __, defaults = inspect.getargspec(v)
                if defaults:
                    print sys.argv[0], k, str(args[:-len(defaults)])[1:-1].replace(",", ""), \
                        str(["%s=%s" % (a, b) for a, b in zip(
                            args[-len(defaults):], defaults)])[1:-1].replace(",", "")
                else:
                    print sys.argv[0], k, str(v.func_code.co_varnames[:v.func_code.co_argcount])[1:-1].replace(",", "")
        sys.exit(-1)
    else:
        func = eval(sys.argv[1])
        args = sys.argv[2:]
        try:
            r = func(*args)
        except Exception, e:
            print "Usage:"
            print "\t", "python %s" % sys.argv[1], str(func.func_code.co_varnames[:func.func_code.co_argcount])[1:-1].replace(",", "")
            if func.func_doc:
                print "\n".join(["\t\t" + line.strip() for line in func.func_doc.strip().split("\n")])
            print e
            r = -1
            import traceback
            traceback.print_exc()
        if isinstance(r, int):
            sys.exit(r)
