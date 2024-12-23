Nebulon distributed / P2P storage supercloud
============================================

Nebulon is a research project designing the building blocks for an encrypted
and highly resilient P2P supercloud.

The basic idea: everybody can just join the network by connecting to some
neighbours, and add storage capacity to it. Data blocks will be replicated
on-demand while traveling through the network.

Key features:
-------------

* on-demand replication
* strong encryption
* inherent deduplication
* plausible deniability
* highly customizable: most operational aspects defined by policies
* pluggable transports

Plausible deniability and censorship resistance
-----------------------------------------------

Since the storage space is distributed over the network of storage nodes,
the network members could easily threatened and politically prosecuted by
oppressive regimes, as it's already the case with other networks, eg. via
the EUSSR's DSA. Therefore having plausible deniability on the foundation
is imperative for survival.

Nebulon is designed to encrypt all data blocks, so individual members have
no way of knowing what's actually stored on their nodes - thus they cannot
be held responsible for that data.

Pluggable architecture
----------------------

The architecture allows for wide range of possible transports and storage
backends, which can be easily implemented as BlockStore drivers.

So far, a REST- and a GRPC-transport have been implemented. Others (e.g. git,
UUCP, SMTP, ...) are planned.

Additional policies, eg. caching and replication, can also be implemented
this way.
