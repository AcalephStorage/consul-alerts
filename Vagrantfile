#!/usr/bin/env ruby
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

default_paths = %w[~/trusty.box ~/trusty64.box trusty.box]

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  
  config.vm.box = "consul-alerts"
  path = default_paths.detect{|path| File.exist?(File.expand_path(path))} || "https://cloud-images.ubuntu.com/vagrant/trusty/current/trusty-server-cloudimg-amd64-vagrant-disk1.box"
  config.vm.box_url = path

  config.vm.provider "virtualbox" do |v|
    v.memory = 500
  end

  config.vm.define 'node1', primary: true do |c|
    c.vm.network "private_network", ip: "192.168.100.50"
    c.vm.hostname = "node1"
  end

  config.vm.define 'node2', primary: true do |c|
    c.vm.network "private_network", ip: "192.168.100.51"
    c.vm.hostname = "node2"
  end

  config.vm.define 'node3', primary: true do |c|
    c.vm.network "private_network", ip: "192.168.100.52"
    c.vm.hostname = "node3"
  end

  config.vm.provision :ansible do |ansible|
    ansible.playbook = 'provision/provision.yml'
  end
end
