var getUrlVars = function() {
	var vars = {};
	var parts = window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function(m,key,value) {
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
};

var zoomOut; // Placeholder for setTimeout callback
function initialize(mapTypeName) {
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
		maxZoom: 1000,
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
		map.setZoom(map.getZoom() + 0.00001);
		zoomOut = function() {
			map.setZoom(map.getZoom() - 0.00001)
		};
		setTimeout("zoomOut()", 1);
	};
	return map;
}
