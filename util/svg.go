package util

// Generate with help of OpenAI Codex

func SVGPathToVector2DArray(path string) (points []*Vector2D) {
	var (
		i        int
		cmd      byte
		hasCmd   bool
		curX     float64
		curY     float64
		startX   float64
		startY   float64
	)

	for i < len(path) {
		c := path[i]
		if isSpace(c) || c == ',' {
			i++
			continue
		}

		if isCommand(c) {
			cmd = c
			hasCmd = true
			i++
			continue
		}

		if !hasCmd {
			i++
			continue
		}

		switch cmd {
		case 'M', 'm':
			x, ok := readNumber(path, &i)
			if !ok {
				break
			}
			y, ok := readNumber(path, &i)
			if !ok {
				break
			}
			if cmd == 'm' {
				curX += x
				curY += y
			} else {
				curX = x
				curY = y
			}
			startX, startY = curX, curY
			points = append(points, Vector(curX, curY))
			cmd = toLine(cmd)
		case 'L', 'l':
			x, ok := readNumber(path, &i)
			if !ok {
				break
			}
			y, ok := readNumber(path, &i)
			if !ok {
				break
			}
			if cmd == 'l' {
				curX += x
				curY += y
			} else {
				curX = x
				curY = y
			}
			points = append(points, Vector(curX, curY))
		case 'H', 'h':
			x, ok := readNumber(path, &i)
			if !ok {
				break
			}
			if cmd == 'h' {
				curX += x
			} else {
				curX = x
			}
			points = append(points, Vector(curX, curY))
		case 'V', 'v':
			y, ok := readNumber(path, &i)
			if !ok {
				break
			}
			if cmd == 'v' {
				curY += y
			} else {
				curY = y
			}
			points = append(points, Vector(curX, curY))
		case 'Z', 'z':
			curX, curY = startX, startY
			i++
		default:
			// Skip unsupported command numbers.
			if _, ok := readNumber(path, &i); !ok {
				i++
			}
		}
	}

	return
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\t' || c == '\r'
}

func isCommand(c byte) bool {
	switch c {
	case 'M', 'm', 'L', 'l', 'H', 'h', 'V', 'v', 'Z', 'z':
		return true
	default:
		return false
	}
}

func toLine(cmd byte) byte {
	if cmd == 'm' {
		return 'l'
	}
	return 'L'
}

func readNumber(s string, idx *int) (float64, bool) {
	i := *idx
	for i < len(s) && (isSpace(s[i]) || s[i] == ',') {
		i++
	}
	if i >= len(s) {
		*idx = i
		return 0, false
	}

	start := i
	if s[i] == '+' || s[i] == '-' {
		i++
	}
	dotSeen := false
	expSeen := false
	for i < len(s) {
		c := s[i]
		if c >= '0' && c <= '9' {
			i++
			continue
		}
		if c == '.' && !dotSeen && !expSeen {
			dotSeen = true
			i++
			continue
		}
		if (c == 'e' || c == 'E') && !expSeen {
			expSeen = true
			i++
			if i < len(s) && (s[i] == '+' || s[i] == '-') {
				i++
			}
			continue
		}
		break
	}

	if start == i {
		*idx = i
		return 0, false
	}

	val, ok := atof(s[start:i])
	*idx = i
	return val, ok
}

func atof(s string) (float64, bool) {
	var (
		i      int
		sign   float64 = 1
		intPart float64
		frac    float64
		denom   float64 = 1
		expSign float64 = 1
		exp     int
	)

	if len(s) == 0 {
		return 0, false
	}
	if s[i] == '-' {
		sign = -1
		i++
	} else if s[i] == '+' {
		i++
	}

	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		intPart = intPart*10 + float64(s[i]-'0')
		i++
	}

	if i < len(s) && s[i] == '.' {
		i++
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			frac = frac*10 + float64(s[i]-'0')
			denom *= 10
			i++
		}
	}

	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		if i < len(s) && s[i] == '-' {
			expSign = -1
			i++
		} else if i < len(s) && s[i] == '+' {
			i++
		}
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			exp = exp*10 + int(s[i]-'0')
			i++
		}
	}

	val := sign * (intPart + frac/denom)
	if exp != 0 {
		if expSign < 0 {
			for j := 0; j < exp; j++ {
				val /= 10
			}
		} else {
			for j := 0; j < exp; j++ {
				val *= 10
			}
		}
	}

	return val, true
}
