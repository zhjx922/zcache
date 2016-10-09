<?php
class Test {
    private $cache;

    public function __construct() {
        $this->cache = new Memcache;
        $this->cache->connect('127.0.0.1', 11211);
    }

    public function set() {
        echo "开始执行set" . PHP_EOL;
        $i = 0;
        while($i++ < 10) {
            var_dump("set", $this->cache->set("k-{$i}", "v-{$i}"));
        }
    }

    public function get() {
        $i = 0;
        while($i++ < 10) {
            var_dump("get", $this->cache->get("k-{$i}"));
        }
    }

    public function delete() {
        $i = 0;
        while($i++ < 10) {
            var_dump("delete", $this->cache->delete("k-{$i}"));
        }
    }

    public function expire() {
        $this->cache->set("expire", serialize(array('a'=>'b')), 1, 4);
        var_dump($this->cache->get("expire"));
        sleep(5);
        var_dump($this->cache->get("expire"));
    }

    public function setLong() {
        $this->cache->set("long", str_repeat("a", 4096));
        $this->cache->set("long2", str_repeat("b", 4096));
    }

    public function getLong() {
        var_dump($this->cache->get("long"));
        var_dump($this->cache->get("long2"));
    }
}

$test = new Test();
$test->set();
$test->get();
$test->delete();
$test->get();
$test->expire();
$test->setLong();
$test->getLong();