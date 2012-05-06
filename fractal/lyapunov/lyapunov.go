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
/*
 * Implements http://en.wikipedia.org/wiki/Lyapunov_fractal
 * For runtime comparison see http://www.efg2.com/Lab/FractalsAndChaos/Lyapunov.htm
 */
package lyapunov

import (
	"fmt"
	"code.google.com/p/go-fracserv/fractal"
	"image"
	"image/color"
	"log"
	"math"
	"strconv"
	"strings"
)

type debugT bool

var debug = debugT(false)

func (d debugT) Println(args ...interface{}) {
	if d {
		log.Println(args...)
	}
}

func (d debugT) Printf(format string, args ...interface{}) {
	if d {
		log.Printf(format, args...)
	}
}

type Lyapunov struct {
	*image.RGBA
	S string
	N int
}

func NewFractal(o fractal.Options) (fractal.Fractal, error) {
	w, err := strconv.Atoi(o.Get("w"))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse width %q: %s", o.Get("w"), err)
	}
	h, err := strconv.Atoi(o.Get("h"))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse height %q: %s", o.Get("h"), err)
	}

	n, err := strconv.Atoi(o.Get("n"))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse num iterations %q: %s", o.Get("n"), err)
	}

	s := strings.ToUpper(o.Get("s"))
	if len(s) == 0 {
		return nil, fmt.Errorf("Failed to parse sequence %q: %s", o.Get("s"), err)
	}

	return &Lyapunov{image.NewRGBA(image.Rect(0, 0, w, h)), s, n}, nil
}

func (l *Lyapunov) At(x, y int) color.Color {
	bounds := l.Bounds()

	a := float64(1 + (x/bounds.Dx())*4)
	b := float64(1 + (y/bounds.Dy())*4)

	x0 := float64(0.5)
	debug.Println("First pass")
	for n := range l.S {
		var r float64
		switch {
		case l.S[n] == 'A':
			r = a
		case l.S[n] == 'B':
			r = b
		default:
			log.Fatalf("Sequence value not A or B: %q", l.S[n])
		}
		x0 = r * x0 * (1 - x0)
		debug.Println("x0", x0, "a", a, "b", b, "r", r)
	}

	// 	double sum_log_deriv = 0;
	// 	for (int n = 0; n < numRounds; n++)
	// 	{
	// 		double prod_deriv = 1;
	// 		for (int m = 0; m < seq_length; m++)
	// 		{
	// 			r = seq[m] == 1 ? b : a;
	// 			/* avoid computing too many logarithms. One every round is acceptable. */
	// 			prod_deriv *= r * (1 - 2 * x);
	// 			x = r * x * (1 - x);
	// 		}
	// 		double deriv_log = Math.Log(Math.Abs(prod_deriv));
	// 		sum_log_deriv += deriv_log;
	// 		//Console.WriteLine("(" + xPos + "," + yPos + ") Iter " + n + " Log " + deriv_log);
	// 	}
	// 	double lambda = sum_log_deriv / (numRounds * seq_length);

	debug.Println("Second pass")
	sumLogDeriv := float64(0)
	for i := 0; i < l.N; i++ {
		prodDeriv := float64(1)
		for n := range l.S {
			var r float64
			switch {
			case l.S[n] == 'A':
				r = a
			case l.S[n] == 'B':
				r = b
			}

			prodDeriv *= r * (1 - 2*x0)
			x0 = r * x0 * (1 - x0)
			debug.Println("x0", x0)
		}
		derivLog := math.Log(math.Abs(prodDeriv))
		sumLogDeriv += derivLog
		debug.Println("derivLog", derivLog)
		debug.Println("sumLogDeriv", sumLogDeriv)
	}
	lambda := sumLogDeriv / float64(l.N*len(l.S))
	col := lambda
	if col < 0 {
		col += 1
	}
	debug.Println("lambda", lambda, "col", col)

	return color.RGBA{uint8(255 * col), 0, 0, 255}
}

func (m *Lyapunov) Ratio() float32 {
	return 1.0
}
