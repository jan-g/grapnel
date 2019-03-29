# Fn reverse shell (grappling hook)

Since it's only typically a matter of time before these appear - AND because they're
pretty useful - it seems pertinent to try this directly.

This set of scripts constructs a small shell environment within a running Fn serverless function.

In order to use it, you'll need a bastion to use as a jump-host.

The idea is simple: the function sets up an outgoing ssh connection to the jump-host.
Once attached, it constructs a reverse tunnel which feeds ssh connections back to an
ssh server running inside the function.

By continually invoking the function, the shell can be kept alive and active.

(This is based upon an idea demoed in a [talk](https://www.meetup.com/bristech/events/251644716/)
for the local ["Bris-Tech"](https://bris.tech) group given in July 2018.)

## Using the grapnel

First, you'll need to set up your bastion host. I'd recommend that you put this on the
same subnet that your functions will appear on, to keep communication between the
function and the jump-host private.

Supposing that you have a host whose default user is `opc` on `10.0.0.2` (here we
*will* indeed be using a private subnet for Fn <-> jumphost communication):

    ./set-up-app -h 10.0.0.2 -u opc
    
Add the generated public key to that host, together with the private key used
to authenticate the connection in the other direction:

    scp f_to_h.pub h_to_f opc@jumphost

Run the loop to keep the function active.

    ./keep-function-running

*On the jumphost*, you should see a listening port open waiting for your connection.
You can then get back into the function container and look around:

    [opc@jumphost ~]$ ssh -v -p 9999 -i h_to_f -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@localhost
    ...
    f9be91b59137:~# ps
    PID   USER     TIME  COMMAND
        1 root      0:00 /dev/init -- /start
        6 root      0:00 ./func
        7 root      0:00 {start} /bin/sh /start
       19 root      0:00 dropbear -R -E -s -p 9999
       26 root      0:00 ssh -v -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -i /tmp/.ssh/id_rsa -R9999:localhost:999
       49 root      0:00 dropbear -R -E -s -p 9999
       50 root      0:00 -ash
       53 root      0:00 ps