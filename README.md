# PiSec-Proxy

##The local network proxy of the PiSec project (https://github.com/ringloop/pisec) 

Our aim in implementing the Proxy, and the whole PiSec project, was to protect our local network from malware, and data integrity.

So, the hero we needed to inspire to is no one less than the mighty *Gandalf*. 

<p align="center">
  <img src="https://user-images.githubusercontent.com/4531376/189697521-43557602-5a73-48dd-8a74-4a8c4daeda0f.png">
</p>

>>thanks to Anja van Hagen for creating this beautiful image, you can check it out at: https://codepen.io/anjanas_dh/pen/ZMqKwb

So as he defeated the Balrog in the underground realm of Moria, the Proxy will provide protection checking 

The proxy is designed to run in the local network of the customer, receive all the HTTP requests and filter them, blocking navigation in case a link is condsidered harmful, and we invited the most powerful blocker of the middle earth to help us in this: the one and only Gandalf 

The proxy downloads the Bloom Filter from the server at startup, for each request received, the bloom filter is checked and if it matches positive, this means that the link is _probably_ harmful. 

Due the uncertainty of Bloom Filter results, in case of positive match, a second confirmation from the server is necessary to make the result sure. 
To keep system performance acceptable, we implemented a local cache as a Redis (www.redis.io) repository. Here we keep 3 sets of data:

- Denied: this index will contains all the confirmed dangerous links. Navigation to these sites should be blocked. 

- False Positives: this index contains all the links giving a positive match in bloom filter, but it has been proven as a false positive by the second match. Navigation to this site should be always granted. 

- Allowed: this index contains all those links that has been granted by user although it has been confirmed as harmful by the system. The navigation to these links will always be granted. 



