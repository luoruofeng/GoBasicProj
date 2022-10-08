# GoBasicProj
GoBasicProj is basic master/slaver(master/worker) mode template of golang project.

The project consists of two parts: Master and Worker.

ETCD as message middleware. ETCD as data transfer proxy.

Master send task to ETCD and worker watching task from ETCD.

