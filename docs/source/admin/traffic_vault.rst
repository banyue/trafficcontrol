..
..
.. Licensed under the Apache License, Version 2.0 (the "License");
.. you may not use this file except in compliance with the License.
.. You may obtain a copy of the License at
..
..     http://www.apache.org/licenses/LICENSE-2.0
..
.. Unless required by applicable law or agreed to in writing, software
.. distributed under the License is distributed on an "AS IS" BASIS,
.. WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
.. See the License for the specific language governing permissions and
.. limitations under the License.
..

****************************
Traffic Vault Administration
****************************
Installing Traffic Vault
========================
In order to successfully store private keys you will need to install Riak.
The latest version of Riak can be downloaded on the Riak `website <http://docs.basho.com/riak/latest/downloads/>`_.
The installation instructions for Riak can be found `here <http://docs.basho.com/riak/latest/ops/building/installing/>`__.

Based on experience, version 2.0.5 of Riak is recommended, but the latest version should suffice.

Configuring Traffic Vault
=========================
The following steps were taken to configure Riak in Comcast production environments.


Self Signed Certificate configuration
-------------------------------------

.. note:: Self Signed Certificates are not recommended for production use. Intended for dev or learning purposes only. Modify subject as necessary.   

.. code-block:: shell
    :caption: Self-Signed Certificate Configuration

	cd ~
	mkdir certs
	cd certs
	openssl genrsa -out ca-bundle.key 2048
	openssl req -new -key ca-bundle.key -out ca-bundle.csr -subj "/C=US/ST=CO/L=DEN/O=somecompany/OU=CDN/CN=somecompany.net/emailAddress=someuser@somecompany.net"
	openssl x509 -req -days 365 -in ca-bundle.csr -signkey ca-bundle.key -out ca-bundle.crt
	openssl genrsa -out server.key 2048
	openssl req -new -key server.key -out server.csr -subj "/C=US/ST=CO/L=DEN/O=somecompany/OU=CDN/CN=somecompany.net/emailAddress=someuser@somecompany.net"
	openssl x509 -req -days 365 -in server.csr -CA ca-bundle.crt -CAkey ca-bundle.key -CAcreateserial -out server.crt
	mkdir /etc/riak/certs
	mv -f server.crt /etc/riak/certs/.
	mv -f server.key /etc/riak/certs/.
	mv -f ca-bundle.crt /etc/pki/tls/certs/.


Riak configuration file configuration
-------------------------------------

The following steps need to be performed on each Riak server in the cluster:

* Log into Riak server as root

* ``cd /etc/riak/``

* Update the following in ``riak.conf`` to reflect your IP, hostname and CDN domains/sub-domains:
	- ``nodename = riak@a-host.sys.kabletown.net``
	- ``listener.http.internal = a-host.sys.kabletown.net:8098`` (port can be 80 - This endpoint will not work with sec enabled)
	- ``listener.protobuf.internal = a-host.sys.kabletown.net:8087`` (can be different port if you want)
	- ``listener.https.internal = a-host.sys.kabletown.net:8088`` (port can be 443)

* Updated the following conf file to point to your cert files
	- ``ssl.certfile = /etc/riak/certs/server.crt``
	- ``ssl.keyfile = /etc/riak/certs/server.key``
	- ``ssl.cacertfile = /etc/pki/tls/certs/ca-bundle.crt``

* Add a line at the bottom of the configuration file for TLSv1
	- ``tls_protocols.tlsv1 = on``

* Once the configuration file has been updated restart Riak
	- ``/etc/init.d/riak restart``

* Validate server is running by going to the following URL:
 	- ``https://<serverHostname>:8088/ping``

``riak-admin`` configuration
----------------------------

``riak-admin`` is a command line utility that needs to be run as root on a server in the Riak cluster.

Assumptions:
	* Riak 2.0.2 or greater is installed
	* SSL Certificates have been generated (signed or self-signed)
	* Root access to Riak servers

Add ``admin`` user and ``riakuser`` to Riak
	* ``admin`` user will be a super user
	* ``riakuser`` will be the application user

Login to one of the riak servers in the cluster as root (any will do)

	1. Enable security

		``riak-admin security enable``

	2. Add groups

		``riak-admin security add-group admins``

		``riak-admin security add-group keysusers``
	3. Add users

		.. Note:: User name and password should be stored in ``/opt/traffic_ops/app/conf/<environment>/riak.conf``

		``riak-admin security add-user admin password=<AdminPassword> groups=admins``

		``riak-admin security add-user riakuser password=<RiakUserPassword> groups=keysusers``

	4. Grant access for ``admin`` and ``riakuser``

		``riak-admin security add-source riakuser 0.0.0.0/0 password``

		``riak-admin security add-source admin 0.0.0.0/0 password``

	5. Grant privileges to the ``admins`` group for everything

		``riak-admin security grant riak_kv.list_buckets,riak_kv.list_keys,riak_kv.get,riak_kv.put,riak_kv.delete on any to admins``

	6. Grant privileges to ``keysusers`` group for SSL, DNSSEC, and ``url_sig_keys`` buckets only

		``riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete on default ssl to keysusers``

		``riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete on default dnssec to keysusers``

		``riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete on default url_sig_keys to keysusers``

		``riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete on default cdn_uri_sig_keys to keysusers``

.. seealso:: For more information on security in Riak, see the `Riak Security documentation <http://docs.basho.com/riak/2.0.4/ops/advanced/security/>`_.
.. seealso:: For more information on authentication and authorization in Riak, see the `Riak Authentication and Authorization documentation <http://docs.basho.com/riak/2.0.4/ops/running/authz/>`_.


Traffic Ops Configuration
-------------------------

There are a couple configurations that are necessary in Traffic Ops.

1. Database Updates
	* The servers in the Riak cluster need to be added to the server table (TCP Port = 8088, type = RIAK, profile = RIAK_ALL)

2. Configuration updates
	* ``/opt/traffic_ops/app/conf/<environment>/riak.conf`` needs to be updated to reflect the correct username and password for accessing riak.

Configuring Riak Search
=======================

In order to more effectively support retrieval of SSL certificates by Traffic Router and Traffic Ops ORT, Traffic Vault uses `Riak search <http://docs.basho.com/riak/kv/latest/using/reference/search/>`_. Riak Search uses `Apache Solr <http://lucene.apache.org/solr>`_ for indexing and searching of records. The following explains how to enable, configure, and validate Riak Search.

Riak Configuration
------------------

On Each Riak Server:

1. If Java (JDKv1.8+) is not already installed on your Riak server, install Java
	* To see if Java is already installed: ``java -version``
	* To install Java: ``yum install -y java-1.8.0-openjdk java-1.8.0-openjdk-devel`` (CentOS/RedHat/Fedora), ``apt-get install -y openjdk-8-jdk`` (Ubuntu/Debian/Linux Mint), ``pacman -Sy jdk8-openjdk`` (Arch/Manjaro)

2. Enable search in riak.conf
	* ``$EDITOR /etc/riak/riak.conf``
	* look for search and change ``search = off`` to ``search = on``

3. Restart Riak so search is on
	* ``systemctl restart riak`` (systemD-based systems)

One time configuration:

1. On one of the Riak servers in the cluster run the following riak-admin commands:

	- ``riak-admin security grant search.admin on schema to admin``

	- ``riak-admin security grant search.admin on index to admin``

	- ``riak-admin security grant search.query on index to admin``

	- ``riak-admin security grant search.query on index sslkeys to admin``

	- ``riak-admin security grant search.query on index to riakuser``

	- ``riak-admin security grant search.query on index sslkeys to riakuser``

	- ``riak-admin security grant riak_core.set_bucket on any to admin``

2. Add the search schema to Riak. This schema is a simple Apache Solr configuration file which will index all records on CDN, hostname, and Delivery Service.
	* Get the schema file by either cloning the project and going to ``traffic_ops/app/config/misc/riak_search`` or from `Github <https://github.com/apache/trafficcontrol/tree/master/traffic_ops/app/conf/misc/riak_search>`_.
	* Use cURL to add the schema to Riak: ``curl -kvs -XPUT "https://admin:pass@riakserver:8088/search/schema/sslkeys" -H 'Content-Type:application/xml' -d @sslkeys.xml``

3. Add search index to Riak
	* run the following cURL command: ``curl -kvs -XPUT "https://admin:pass@riakserver:8088/search/index/sslkeys" -H 'Content-Type: application/json' -d '{"schema":"sslkeys"}'``

4. Associate the ``sslkeys`` index to the ``ssl`` bucket in Riak
	* run the following curl command: ``curl -kvs -XPUT "https://admin:pass@riakserver:8088/buckets/ssl/props" -H'content-type:application/json' -d'{"props":{"search_index":"sslkeys"}}'``

Riak Search (using Apache Solr) will now index all **new** records that are added to the ``ssl`` bucket. The ``cdn``, ``deliveryservice``, and ``hostname`` fields are indexed. When a search is performed Riak will return the indexed fields along with the certificate and key values for a SSL record. In order to add the indexed fields to current records and to get the current records added, a standalone script needs to be run. This does not need to be done on new installs. The following explains how to run the script.

1. Get script from Github either by cloning the project and going to ``traffic_ops/app/script`` or from `here <https://github.com/apache/trafficcontrol/blob/master/traffic_ops/app/script/update_riak_for_search.pl>`_
2. Run the script by performing the following command ``./update_riak_for_search.pl -to_url=https://traffic-ops.kabletown.net -to_un=user -to_pw=password`` (with the appropriate URL substituted for your Traffic Ops server{.})

Validate the search is working by querying against Riak directly:
``curl -kvs "https://admin:password@riakserver:8088/search/query/sslkeys?wt=json&q=cdn:mycdn"``

Validation can also be done by querying Traffic Ops:
``curl -Lvs -H "Cookie: $COOKIE" https://traffic-ops.kabletown.net/api/1.2/cdns/name/mycdn/sslkeys.json``
