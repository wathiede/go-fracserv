// Copyright 2012 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
var getUrlVars = function (url) {
	var vars = {};
	var url;
	if (url == undefined) {
		url = window.location.href;
	}
	var parts = url.replace(/[?&]+([^=&]+)=([^&]*)/gi, function (m, key, value) {
		vars[key] = value;
	});
	return vars;
}

var size = 128;
var serializeArray = function (arr) {
	var out = new Array();

	for (key in arr) {
		out.push(key + '=' + arr[key]);
	}

	return out.join('&');
};

var zoomOut; // Placeholder for setTimeout callback
function initialize(mapTypeName) {
	// TODO(wathiede): set trackResize = true;
	map = L.map('maps').setView([51.505, -0.09], 13);
	L.tileLayer('/debug?w=128&h=128&x={x}&y={y}&z={z}&name=&url=', {
		tileSize: 128,
	}).addTo(map);
	return map;

	var customTypeOptions = {
		getTileUrl: function (coord, zoom) {
			options = {
				w: size,
				h: size,
				x: coord.x,
				y: coord.y,
				z: Math.floor(zoom),
			}
			// Add any form elements to request
			$('form input').each(function (idx, e) {
				options[e.id] = e.value;
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
	map.reload = function () {
		// Crappy hack to make redraw work
		map.setZoom(map.getZoom() + 0.00001);
		zoomOut = function () {
			map.setZoom(map.getZoom() - 0.00001)
		};
		setTimeout("zoomOut()", 1);
	};
	map.fracSave = function () {
		options = {
			c: map.getCenter().toUrlValue(),
			z: map.getZoom(),
		};
		// Add any form elements to request
		$('#config form input').each(function (idx, e) {
			options[e.id] = e.value;
		});
		return "#" + mapTypeName + "?" + serializeArray(options);
	};

	return map;
}
