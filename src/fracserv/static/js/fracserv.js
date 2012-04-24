!function( $ ) {
	$(function () {
		"use strict"
		var fracType = 'solid'
		var getContents = function() {
			var el = $(this);
			$.get('/' + fracType, function(data) {
				el.attr('data-content', data);
				el.attr('title', fracType);
				el.popover('show');
				$('form').submit(function() {
					$('body').css('background-image', 'url(/' + fracType + '?' + $(this).serialize() + ')');
					return false;
				});
				return data;
			})

		};

		$('#config').click(getContents);

		$('#config').popover({
			placement: 'top',
			trigger: 'manual',
		});
	})
}( window.jQuery );
