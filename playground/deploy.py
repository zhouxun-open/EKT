#!/usr/bin/env python
# -*- coding:utf-8 -*-

import os
import json
import glob

from fabric import Connection, SerialGroup

EKTCLI = './ecli'
assert os.path.exists(EKTCLI), 'ecli not found'

HOST_ADDR = [
    '192.168.6.54',
    '192.168.6.55',
    '192.168.6.56',
    '192.168.6.57',
    '192.168.6.58',
    '192.168.6.59'
]
NODE_NUM = 3
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
    for _ in range(num):
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

def _get_account(num):
    cmd = '%s account new' % EKTCLI
    ret = []
    for _ in range(num):
        data = os.popen(cmd).read().strip()
        lines = data.split('\n')
        if len(lines) != 2:
            continue
        items = lines[0].split(':')
        if len(items) != 2 or 'Private Key' not in items[0]:
            continue
        pk = items[1].strip()
        items = lines[1].split(':')
        if len(items) != 2 or 'Your address is' not in items[0]:
            continue
        addr = items[1].strip()
        ret.append((pk, addr))
    return ret

def gen_conf():
    meta = open('peer_info.txt', 'wb')

    num = len(HOST_ADDR) * NODE_NUM
    peers = _gen_peers(num)
    genesis_tpl = open('genesis.tpl').read()
    genesis_tpl = genesis_tpl.replace('{{.env}}', ENV)
    genesis_act = _get_account(3)
    for i, act in enumerate(genesis_act):
        genesis_tpl = genesis_tpl.replace('{{.genesisAddr' + str(i) + '}}', act[1])
        meta.write('pk:%s addr:%s\n' % (act[0], act[1]))
        meta.flush()

    nets = []
    for addr in HOST_ADDR:
        for i in range(NODE_NUM):
            peer = peers.pop(0)
            pk, peerId = peer
            port = PORT_RANGE[i]
            my_conf = genesis_tpl
            my_conf = my_conf.replace('{{.addr}}', addr)
            my_conf = my_conf.replace('{{.port}}', str(port))
            my_conf = my_conf.replace('{{.privateKey}}', pk)
            my_conf = my_conf.replace('{{.peerId}}', peerId)
            my_conf = my_conf.replace('{{.addrVer}}', str(ADDR_VERSION))
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

def publish_conf():
    netenv_conf = 'conf/%s.json' % ENV
    for addr in HOST_ADDR:
        with Connection(addr) as conn:
            conn.run('mkdir -p %s' % REMOTE_CONF_DIR)
            conn.put(netenv_conf, os.path.join(REMOTE_CONF_DIR, os.path.basename(netenv_conf)))
            conn.put('ctrl.py', os.path.join(REMOTE_CONF_DIR, 'ctrl.py'))
            genesis_confs = glob.glob('conf/%s_*.genesis.json' % addr)
            for conf in genesis_confs:
                conn.put(conf, os.path.join(REMOTE_CONF_DIR, os.path.basename(conf)))
    debug('publish conf done.')
        
def deploy():
    publish_conf()
    SerialGroup(*HOST_ADDR).run('cd %s && python ctrl.py pull %s' % (REMOTE_CONF_DIR, IMAGE_URL))
    debug('pull image done.')
    SerialGroup(*HOST_ADDR).run('cd %s && python ctrl.py run %s' % (REMOTE_CONF_DIR, IMAGE_URL))
    debug('run container done.')

# stop ekt8 container
def stop():
    SerialGroup(*HOST_ADDR).run('cd %s && python ctrl.py stop' % REMOTE_CONF_DIR)

# start ekt8 container
def start():
    SerialGroup(*HOST_ADDR).run('cd %s && python ctrl.py start' % REMOTE_CONF_DIR)

# restart ekt8 container
def restart():
    SerialGroup(*HOST_ADDR).run('cd %s && python ctrl.py restart' % REMOTE_CONF_DIR)

# clean ekt8 container
def clean():
    SerialGroup(*HOST_ADDR).run('cd %s && python ctrl.py clean' % REMOTE_CONF_DIR)

# show ekt8 container date
def date():
    SerialGroup(*HOST_ADDR).run('cd %s && python ctrl.py date' % REMOTE_CONF_DIR)

# exec cmd in ekt8 container
def exec_cmd(cmd):
    SerialGroup(*HOST_ADDR).run('cd %s && python ctrl.py exec_cmd "%s"' % (REMOTE_CONF_DIR, cmd))

def upload(src, dst):
    for addr in HOST_ADDR:
        with Connection(addr) as conn:
            conn.put(src, dst)

def run_cmd(cmd):
    for addr in HOST_ADDR:
        with Connection(addr) as conn:
            conn.run(cmd)




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
