视野服务 Area Of Interest (AOI)
=============================================

游戏服务器的AOI（area of interest)部分，位置有关的游戏实体一般都有一个视野或关心的范围

* 适用于大批量的对象管理和查询  
* 采用四叉树来进行管理，可以通过少量改动支持3D  
* 动态节点回收管理，减少内存, 节点的数量不会随着四叉树的层次变得不可控，与对象数量成线性关系  
* 缓存搜索上下文，减少无效搜索，对于零变更的区域，搜索会快速返回  
* 整个搜索过程，支持自定义过滤器进行单元过滤  
* AOI 支持对象的内存缓冲区管理  
* 线程不安全  
* 支持单位半径，不同单位可以定义不同的半径(开启半径感知，会损失一部分性能)  
* 只提供了搜索结构，不同单位的视野区域可以定义不一样  


The aoi module include a set of aoi interface, and an implementation of tower aoi algorithm.  
##Installation
```
npm install pomelo-aoi
```
##Generate an aoi instance
For the aoi service can be used in many areas, each area use the aoi module should use it's own aoi instance.
We use a aoi factory to generate aoi instance, it accept an object as parameter, and return an aoi instance,  which can be used to implament the aoi function.   

``` javascript
var aoiManager = require('pomelo-aoi');
var config = {
    map : {
        width : 3200,
        height : 2400
    },
    tower : {
        width : 300,
        height : 300
    }
}

var aoi = qoiManager.getService(config);
```

##Use the aoi service
The aoi instace has the basic interface for aoi action.

``` javascript
    //Add object 
    aoi.addObject(obj, pos);
    
    //Remove object 
    aoi.removeObject(obj, pos);
    
    //Update object
    aoi.updateObject(obj, oldPos, newPos);
    
    //Add watcher 
    aoi.addWatcher(watcher, pos, range);
    
    //Remove watcher
    aoi.removeWatcher(watcher, pos, range0;
    
    //updateWatcher(watcher, oldPos, newPos, oldRange, newRange);
``` 
More api can be find in aoiService.js.

##Handle aoi event
The aoi service will generate event when the status of objects or watchers changes. You can handler these event :
``` javascript
    aoi.on('add', function(params){
        //Handle add event
    });

``` 
The event of tower aoi are: 'add', 'remove', 'update' for aoi object, and 'updateWatcher' for watcher.
Of course you can ignore all these events without do any effect to aoi function. 
