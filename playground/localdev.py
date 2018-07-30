#!/usr/bin/env python
# -*- coding:utf-8 -*-

import os
import json
import glob

DOCKER0_ADDR = '172.17.0.1'
EKTCLI = './ecli'
assert os.path.exists(EKTCLI), 'ecli not found'
NODE_NUM = 3
LOCALDEV_PORT = 19993
PORT_RANGE = [19990, 19991, 19992]
assert len(PORT_RANGE) >= NODE_NUM, "port not enough"

ENV = 'localnet'
ADDR_VERSION = 4
IMAGE_URL = 'registry.cloudhua.com/ekt8/ekt8:latest'
LOCAL_CONF_DIR = 'conf'
REMOTE_CONF_DIR = '/data/EKT/conf'

os.system('mkdir -p %s' % LOCAL_CONF_DIR)

DEBUG = True
def debug(msg):
    if DEBUG:
        print msg


def _gen_peers(num):
    cmd = '%s node init' % EKTCLI
    ret = [] 
    for _ in range(num+1):
        data = os.popen(cmd).read().strip()
        lines = data.split('\n')
        if len(lines) != 2:
            continue
        items = lines[0].split(':')
        if len(items) != 2 or 'private key' not in items[0]:
            continue
        pk = items[1].strip()
        items = lines[1].split(':')
        if len(items) != 2 or 'peerId' not in items[0]:
            continue
        peerId = items[1].strip()
        ret.append((pk, peerId))
    return ret

def gen_conf():
    num = NODE_NUM
    peers = _gen_peers(num)
    genesis_tpl = open('genesis.tpl').read()
    genesis_tpl = genesis_tpl.replace('{{.env}}', ENV)
    nets = []
    meta = open('peer_info.txt', 'wb')
    addr = DOCKER0_ADDR
    PORT_RANGE.append(LOCALDEV_PORT)
    for i in range(NODE_NUM+1):
        peer = peers.pop(0)
        pk, peerId = peer
        port = PORT_RANGE[i]
        my_conf = genesis_tpl
        my_conf = my_conf.replace('{{.addr}}', addr)
        my_conf = my_conf.replace('{{.port}}', str(port))
        my_conf = my_conf.replace('{{.privateKey}}', pk)
        my_conf = my_conf.replace('{{.peerId}}', peerId)
        my_conf = my_conf.replace('{{.addrVer}}', str(ADDR_VERSION))
        if i == NODE_NUM:
            name = 'genesis.localdev.json'
        else:
            name = '%s_%s.genesis.json' % (addr, port)
        save_path = os.path.join(LOCAL_CONF_DIR, name)
        open(save_path, 'wb').write(my_conf)
        debug('gen %s' % save_path)
        nets.append([peerId, addr, port, ADDR_VERSION, ''])
        meta.write('%s %s %s %s\n' % (addr, port, pk, peerId))
        meta.flush()
    save_path = os.path.join(LOCAL_CONF_DIR, '%s.json' % ENV)
    open(save_path, 'wb').write(json.dumps(nets, indent=2))
    meta.close()
    debug('gen %s' % save_path)
    debug('gen conf done.')

def deploy():
    os.system('cp ctrl.py %s' % LOCAL_CONF_DIR)
    os.system('cd %s && python ctrl.py pull %s' % (LOCAL_CONF_DIR, IMAGE_URL))
    debug('pull image done.')
    os.system('cd %s && python ctrl.py run %s' % (LOCAL_CONF_DIR, IMAGE_URL))
    debug('run container done.')

# stop ekt8 container
def stop():
    os.system('cd %s && python ctrl.py stop' % LOCAL_CONF_DIR)

# start ekt8 container
def start():
    os.system('cd %s && python ctrl.py start' % LOCAL_CONF_DIR)

# restart ekt8 container
def restart():
    os.system('cd %s && python ctrl.py restart' % LOCAL_CONF_DIR)

# clean ekt8 container
def clean():
    os.system('cd %s && python ctrl.py clean' % LOCAL_CONF_DIR)

# show ekt8 container date
def date():
    os.system('cd %s && python ctrl.py date' % LOCAL_CONF_DIR)


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
