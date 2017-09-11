# -*- mode: ruby -*-
# vi: set ft=ruby :


Vagrant.configure(2) do |config|

  config.ssh.insert_key = false
  config.vm.box = "ubuntu/trusty64"

  3.times do |i|
    hostname = "consul-#{i}"
    config.vm.define hostname do |node|
      node.vm.hostname = hostname
      node.vm.network :private_network, ip: "192.168.10.#{10+i}"
    end
  end

  config.vm.provision :ansible do |ansible|
    ansible.playbook = 'provision/site.yml'
    ansible.extra_vars = {
      bootstrap_expect: 3
    }
  end

end
