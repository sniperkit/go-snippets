package utils


const M int = 100
const N int = 100
const NALIGN int = 200

const STOP int = 0
const UP int = 1
const LEFT int = 2
const DIAG int = 3

const MATCHSCORE int = 1
const MISMATCHSCORE int = -1
const GAPSCORE int = -1


/*****************************************************************************************************
 * Function: SmithWaterman
 * Purpose: To calculate similarity percentage between two texts using Smith-Waterman algorithm
 * Parameter: array of firstText as first text, array of secondText as second text
 * Output: percentage of similarity
 *****************************************************************************************************/
func SmithWaterman(firstText []string, secondText []string) float32 {

	var i,j,tmp,length int;
	var distance [M][N]int;  	/* distance label matrix */
	var trace [M][N]int;     	/* trace matrix */
	var alignX [NALIGN]string; 	/* aligned X sequence */
	var alignY [NALIGN]string; 	/* aligned Y sequence */
	
	// initialization
    for i:=0; i<=len(firstText); i++ { distance[i][0]=0 } /* do this for i=0,1,...,m */
    for j:=0; j<=len(secondText); j++ { distance[0][j]=0 }
    for i:=0; i<=len(firstText); i++ { trace[i][0]=STOP }
    for j:=0; j<=len(secondText); j++ { trace[0][j]=STOP }
    minDist, minI, minJ := 0, 0, 0
	
	// labeling of all nodes, this is the main loop of the algorithm
    for i:=1; i<=len(firstText); i++ {    /* note: we begin at i=1 ! */
		for j:=1; j<=len(secondText); j++ {

			dist := 0; /* distance to node (i,j) from virtual start node */
			trace[i][j] = STOP;

			if firstText[i-1] == secondText[j-1] {
				tmp = distance[i-1][j-1] - MATCHSCORE
			} else {
				tmp = distance[i-1][j-1] - MISMATCHSCORE
			}
			if tmp < dist {
				dist = tmp
				trace[i][j] = DIAG
			}

			tmp = distance[i-1][j] - GAPSCORE
			if tmp < dist {
				dist = tmp
				trace[i][j] = UP
			}

			tmp = distance[i][j-1] - GAPSCORE
			if tmp < dist {
				dist = tmp
				trace[i][j] = LEFT
			}

			distance[i][j] = dist;

			if dist < minDist { /* keep track of where the minimum score is */
				minDist = dist
				minI = i
				minJ = j
			}
		}
    }
	
	// now create aligned sequences alignY and alignX 
    // note: these are created in reverse order! 
	
    iAlign := 0
	
    // unaligned ends
	
	for i:=len(firstText); i>minI; i-- {
		alignY[iAlign] = "*"
		alignX[iAlign] = firstText[i-1]
		iAlign++
	}
	
	for j:=len(secondText); j>minJ; j-- {
		alignY[iAlign] = secondText[j-1]
		alignX[iAlign] = "*"
		iAlign++
    }
	
	// when we come here we know that i==minI and j==minJ
    // it is from this position we make the jump to the virtual stop node
    
    for trace[i][j] != STOP {
		switch trace[i][j] {
			case DIAG:
				alignY[iAlign] = secondText[j-1]
				alignX[iAlign] = firstText[i-1]
				i--
				j--
				iAlign++
			case LEFT:
				alignY[iAlign] = secondText[j-1] 
				alignX[iAlign] = "-"
				j--
				iAlign++
			case UP:
				alignY[iAlign] = "-"
				alignX[iAlign] = firstText[i-1]
				i--
				iAlign++
		}
    }
	
    // unaligned beginning

    for i > 0 {
		alignY[iAlign] = "*"
		alignX[iAlign] = firstText[i-1]
		i--
		iAlign++
    }

    for j > 0 {
		alignY[iAlign] = secondText[j-1]
		alignX[iAlign] = "*"
		j--
		iAlign++
    }
	
	// count the length as maximum divider of percentage
	if len(firstText) < len(secondText) {
		length = len(firstText)
	} else {
		length = len(secondText)
	}
	
	// return the percentage of similarity
	return (-1* float32(minDist)/float32(length)) * 100
}
