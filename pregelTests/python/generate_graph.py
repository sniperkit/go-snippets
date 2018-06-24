"""Given num vertices and max edges per vertex, produces
list of vertices in the form "from to"

Can be used as input to a graph from our Pregel implementation.


"""
import sys
import random

def main():
    if len(sys.argv) != 3:
        print("Usage: {py} num_vertices max_edges_per_vertex"
                .format(py=sys.argv[0]))
        sys.exit(-1)

    num_vertices = int(sys.argv[1])
    max_edges = int(sys.argv[2])

    all_vertices = list(range(0, num_vertices))

    for v in all_vertices:
        num_edges = random.randrange(max_edges)
        for e in random.sample(all_vertices, num_edges):
            print("{frm} {to}".format(frm=v, to=e))

if __name__ == "__main__":
    main()
