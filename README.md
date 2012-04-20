go.fracserv
===========

Homework assignment for the greatest group ever

Requirements
------------

* Package `main`
  * should run a webserver
  * should take a commandline flag to set the port to listen on
  * should print the URL of homepage, i.e. http://computer.example.com:9999/

* Package `fractal`
  * Should define an image type that can implement fractals by defining At()
  * Subpackage `solid` which renders a fractal using only a specific color (si
  * Create at least one other package that implements a famous fractal
  * Write at least one unit test to verify your `func At(x, y, int) color.Color`
  * Subpackages to register themselves with `fractal` like the image encoders do

* Homepage
  * Links to jump to each fractal type you implement
  * Each fractals page should have a form to tweak the coefficents for that particular fractal
  * Form can GET instead of POST to make embedding the images easier
  * Should provide download links for JPEG and PNG version of the image
  * This is an exercise in writing Go templates, not testing your Web 2.0 skillz, so if you hit a wall with the HTML, ask, don't fret about it.


Bonus Points
------------
* Using `pprof` what are your three most expensive functions when rendering the fractal
* Use multiple go routines to render individual parts of the fractal at the same time
* Set the go runtime to use the number the same number of threads as your computer has cores at runtime automatically


Modules Used
------------
  [0]: http://golang.org/pkg/flag flag
  [1]: http://golang.org/pkg/image image
  [2]: http://golang.org/pkg/image/jpeg image/jpeg
  [3]: http://golang.org/pkg/image/png image/png
  [4]: http://golang.org/pkg/net/http net/http
  [5]: http://golang.org/pkg/html/template html/template
  [6]: http://golang.org/pkg/net/http/pprof net/http/pprof
  [7]: http://golang.org/pkg/net/url net/url
  [8]: http://golang.org/pkg/runtime runtime
  [9]: http://golang.org/pkg/runtime/pprof runtime/pprof
  [10]: http://golang.org/pkg/testing testing
