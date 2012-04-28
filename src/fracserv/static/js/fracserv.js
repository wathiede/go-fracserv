String.prototype.capitalize = function() {
	return this.charAt(0).toUpperCase() + this.slice(1);
}

!function( $ ) {
	$(function () {
		"use strict"
		/*
		var pop = $('#config');
		var getContents = function(fracType) {
			$.get('/' + fracType, function(data) {
				pop.attr('data-content', data);
				pop.attr('title', fracType.capitalize());

				pop.popover('show');

				$('form #w').val($(window).width());
				$('form #h').val($(window).height());


				$('form').submit(function() {
					$('body').css('background-image', 'url(/' + fracType + '?' + $(this).serialize() + ')');
					return false;
				}).submit();
				return data;
			});
		};

		pop.popover({
			placement: 'top',
			trigger: 'manual',
		});

		pop.click(function() {
			pop.popover('toggle');
		});
		*/
		var getContents = function(fracType) {
			$('#config').load('/' + fracType + ' form', function() {
				var form = $('#config form');

				$('#w', form).val($(window).width());
				$('#h', form).val($(window).height());
				form.submit(function() {
					console.log("Form submitted");
					$('body').css('background-image',
						'url(/' + fracType + '?' + form.serialize() + ')');
						return false;
						}).submit();
					});
		};

		$('ul.nav li a').click(function(e) {
			var fracType = $(this).attr('id');
			console.log(fracType);
			getContents(fracType);
			$('#masthead').fadeOut();
			$('#config').fadeIn();
			return false;
		});
	})
}( window.jQuery );
