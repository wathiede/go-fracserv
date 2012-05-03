var getUrlVars = function(url) {
	var vars = {};
	var url;
	if(url == undefined) {
		url = window.location.href;
	}
	var parts = url.replace(/[?&]+([^=&]+)=([^&]*)/gi, function(m,key,value) {
		vars[key] = value;
	});
	return vars;
}

var size = 512;
var serializeArray = function(arr) {
	var out = new Array();

	for (key in arr) {
		out.push(key + '=' + arr[key]);
	}

	return out.join('&');
};


function NoProjection() {
	this.origin = new google.maps.Point(0, 0);
	this.pixelPerDegree = size/360;
}

NoProjection.prototype.fromLatLngToPoint = function(latLng) {
	console.log("NoProjection.prototype.fromLatLngToPoint");
	var origin = this.origin;
	var pixelPerDegree = this.pixelPerDegree;
	return new google.maps.Point(origin.x + latLng.lng() * pixelPerDegree,
			origin.y + latLng.lat() * pixelPerDegree);
};

NoProjection.prototype.fromPointToLatLng = function(point, noWrap) {
	console.log("NoProjection.prototype.fromPointToLatLng");
	var origin = this.origin;
	var pixelPerDegree = this.pixelPerDegree;
	return new google.maps.LatLng((point.y - origin.y)/pixelPerDegree,
			(point.x - origin.x)/pixelPerDegree, noWrap);
};

function FractalMapType(options) {
	this.options = options;
}

FractalMapType.prototype.tileSize = new google.maps.Size(size, size);
FractalMapType.prototype.maxZoom = 21; // TODO make this unlimited, or at least larger
FractalMapType.prototype.getTile = function(tileCoord, zoom, ownerDocument) {
	var div = ownerDocument.createElement('div');
	div.style.width = this.tileSize.width + 'px';
	div.style.height = this.tileSize.height + 'px';
	div.style.backgroundImage = "url('" + this.options.getTileUrl(tileCoord, zoom) + "')";
	return div;
};

/*
FractalMapType.prototype.releaseTile = function(tile) {
	//console.log("releaseTile", tile);
};
*/
FractalMapType.prototype.name = "Fractal";
FractalMapType.prototype.alt = "Fractal tile map type";

var zoomOut; // Placeholder for setTimeout callback
function initialize(mapTypeName) {
	var getTileUrl = function(tileCoord, zoom, ownerDocument) {
		z = Math.floor(zoom);
		options = {
			w: size,
			h: size,
			x: tileCoord.x/(1 << z),
			y: tileCoord.y/(1 << z),
			z: z,
		}
		// Add any form elements to request
		$('form input').each(function(idx, e) {
			options[e.id] =  e.value;
		});
		return "/" + mapTypeName + "?" + serializeArray(options);
	};
	var fractalMapType = new FractalMapType({
		getTileUrl: getTileUrl,
	});
	var myLatlng = new google.maps.LatLng(0, 0);
	var myOptions = {
		center: myLatlng,
		zoom: 0,
		streetViewControl: false,
		mapTypeControlOptions: {
			mapTypeIds: []
		},
		zoomControlOptions: {
			style: google.maps.ZoomControlStyle.SMALL
		}
	};

	var map = new google.maps.Map(document.getElementById("maps"), myOptions);
	map.mapTypes.set(mapTypeName, fractalMapType);
	map.setMapTypeId(mapTypeName);
	map.panTo(myLatlng);
	map.reload = function() {
		// Crappy hack to make redraw work
		map.setZoom(map.getZoom() + 0.00001);
		zoomOut = function() {
			map.setZoom(map.getZoom() - 0.00001)
		};
		setTimeout("zoomOut()", 1);
	};
	map.fracSave = function() {
		options = {
			c: map.getCenter().toUrlValue(),
			z: map.getZoom(),
		};
		// Add any form elements to request
		$('form input').each(function(idx, e) {
			options[e.id] =  e.value;
		});
		return "#" + mapTypeName + "?" + serializeArray(options);
	};
	map.fromLatLngToTileCoord = function(latlng) {
		var numTiles = 1 << this.getZoom();
		var projection = this.projection;
		console.log("latlng", latlng.lat(), latlng.lng());
		var worldCoordinate = projection.fromLatLngToPoint(latlng);
		console.log("worldCoordinate", worldCoordinate.x, worldCoordinate.y);
		var pixelCoordinate = new google.maps.Point(
				worldCoordinate.x * numTiles,
				worldCoordinate.y * numTiles);
		console.log("pixelCoordinate", pixelCoordinate.x, pixelCoordinate.y);
		var tileCoordinate = new google.maps.Point(
				Math.floor(pixelCoordinate.x / size),
				Math.floor(pixelCoordinate.y / size));
		console.log("tileCoordinate", tileCoordinate.x, tileCoordinate.y);
		return tileCoordinate;

	};
	map.projectionInfo = function() {
		var bb = map.getBounds();
		var tne = this.fromLatLngToTileCoord(bb.getNorthEast());
		var tsw = this.fromLatLngToTileCoord(bb.getSouthWest());
		console.log("tsw.x", tsw.x, "tsw.y", tsw.y);
		console.log("tne.x", tne.x, "tne.y", tne.y);
	};
	map.projection = new NoProjection();

	return map;
}
