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

var size = 128;
var serializeArray = function(arr) {
	var out = new Array();

	for (key in arr) {
		out.push(key + '=' + arr[key]);
	}

	return out.join('&');
}

function NoProjection() {
}

NoProjection.prototype.fromLatLngToPoint(latlng) {
	return new google.maps.Point(latlng.lat(), latlng.lng());
}

NoProjection.prototype.fromPointToLatLng(point) {
	return new google.maps.LatLng(point.x, point.y);
}

var zoomOut; // Placeholder for setTimeout callback
function initialize(mapTypeName) {
	var fromLatLngToTileCoord = function(latlng) {
		var numTiles = 1 << this.getZoom();
		var projection = this.getProjection();
		var worldCoordinate = projection.fromLatLngToPoint(latlng);
		var pixelCoordinate = new google.maps.Point(
				worldCoordinate.x * numTiles,
				worldCoordinate.y * numTiles);
		return new google.maps.Point(
				Math.floor(pixelCoordinate.x / size),
				Math.floor(pixelCoordinate.y / size));

	};
	var projectionInfo = function() {
		var bb = map.getBounds();
		var tne = this.fromLatLngToTileCoord(bb.getNorthEast());
		var tsw = this.fromLatLngToTileCoord(bb.getSouthWest());
		console.log("tsw.x", tsw.x, "tsw.y", tsw.y);
		console.log("tne.x", tne.x, "tne.y", tne.y);
	};
	var customTypeOptions = {
		getTileUrl: function(coord, zoom) {
			options = {
				w: size,
				h: size,
				x: coord.x,
				y: coord.y,
				z: Math.floor(zoom),
			}
			// Add any form elements to request
			$('form input').each(function(idx, e) {
				options[e.id] =  e.value;
			});
			return "/" + mapTypeName + "?" + serializeArray(options);
		},
		tileSize: new google.maps.Size(size, size),
		maxZoom: 21, // TODO make this unlimited, or at least larger
		minZoom: 0,
		name: mapTypeName
	};
	var customMapType = new google.maps.ImageMapType(customTypeOptions);
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
	map.mapTypes.set(mapTypeName, customMapType);
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
	map.projectionInfo = projectionInfo;
	map.fromLatLngToTileCoord = fromLatLngToTileCoord;
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

	return map;
}
