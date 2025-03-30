

/*	Apply assetKeys script by eishiya, last updated Sep 28 2023

	Adds an action to the Map menu that sets "assetKey" on all Tile Objects
	to the tile's file name.
*/

var applyTileObjectProperties = tiled.registerAction("applyAssetKeys", function(action) {
    var map = tiled.activeAsset;
    if(!map || !map.isTileMap)
        return;

    function copyProperties(curLayer) {
        if(curLayer.isObjectLayer) {
            let objects = curLayer.objects;
            for(let obj = 0; obj < objects.length; ++obj) {
                let object = objects[obj];
                if(object.tile) {
                    let tile = object.tile;
                    let tileName = tile.imageFileName;
                    if(!tileName || tileName.length == 0) tileName = tile.tileset.image; //fallback to tilesheet name for non-Collections
                    if(tileName && tileName.length > 0) {
                        tileName = FileInfo.fileName(tileName); // path/to/tile.png -> tile.png
                        object.setProperty("assetKey", tileName);
                    }
                }
            }
        } else if(curLayer.isGroupLayer || curLayer.isTileMap) {
            var numLayers = curLayer.layerCount;
            for(var layerID = 0; layerID < numLayers; layerID++) {
                copyProperties(curLayer.layerAt(layerID));
            }
        }
    }

    copyProperties(map);
});

applyTileObjectProperties.text = "Set assetKeys";
tiled.extendMenu("Map", [
    { action: "applyAssetKeys", before: "MapProperties" }
]);

//Run this action automatically when saving the map:
tiled.assetAboutToBeSaved.connect(function(asset) {
    if(!asset.isTileMap) return;
    let prevAsset = tiled.activeAsset; //save the user's currently viewed asset
    tiled.activeAsset = asset;
    tiled.trigger("applyAssetKeys");
    tiled.activeAsset = prevAsset; //bring back user's viewed asset
});
