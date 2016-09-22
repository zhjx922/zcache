<?php
class Test {
    private $cache;

    public function __construct() {
        $this->cache = new Memcache;
        $this->cache->connect('127.0.0.1', 11211);
    }

    public function set() {
        $i = 0;
        while($i++ < 100000) {
            var_dump($this->cache->set("k-{$i}", "v-{$i}", 0));
        }
    }

    public function get() {
        $i = 0;
        while($i++ < 100000) {
            var_dump($this->cache->get("k-{$i}"));
        }
    }

    public function delete() {
        $i = 0;
        while($i++ < 100000) {
            var_dump($this->cache->delete("k-{$i}"));
        }
    }
}

$test = new Test();
$test->set();
$test->get();
$test->delete();
$test->get();