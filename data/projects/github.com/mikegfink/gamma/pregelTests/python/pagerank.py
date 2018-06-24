"""This is an adaptation of the sequential, matrix-based pagerank
algorithm presented in as pagerank_test in
http://www.michaelnielsen.org/ddi/pregel/

To use this:
    python3 testpagerank.py testfilename

    Produces list of vertex id and pagerank
    e.g.,
    0 0.0588083224524
    1 0.102275343395
    2 0.0653001502556
    3 0.108767171199
    4 0.0485017092391

testfilename is a space-separated list of from-to edges. E.g.,
0   2
0   3
1   1

10  2
# this is a comment; these and blank lines are ignored

Assumption: vertex ids are consecutive ints starting at 0
Any id smaller than the max id, but that is never mentioned testfile is
assumed to represent a vertex with no outgoing or incoming edges.

"""

from numpy import *
import sys
import fileinput

class Vertex():

    def __init__(self, id):
        # This is mostly self-explanatory, but has a few quirks:
        self.id = id
        self.value = 1
        self.out_vertices = [] # [Vertex]

    def add_edge(self, vertex):
        self.out_vertices.append(vertex)

    def __repr__(self):
        outs = " ".join([str(v.id) for v in self.out_vertices])

        return "{vid}->({outs})".format(vid=self.id, outs=outs)

def main():
    if len(sys.argv) > 2:
        print("Usage: {py} test_filename".format(py=sys.argv[0]))
        sys.exit(-1)

    vertices = load_graph()
    num_vertices = max([v.id for v in vertices]) + 1
    set_init_value(vertices, num_vertices)
    page_ranks = pagerank_test(vertices, num_vertices)
    print_prs(page_ranks)

def load_graph():
    vertices = {}

    for line in fileinput.input():
        line = line.strip()
        if len(line) == 0 or line[0] == '#':
            continue
        e = line.split()
        from_id = int(e[0]) -1
        to_id = int(e[1]) -1
        edge(vertices, from_id, to_id)

    return list(vertices.values())

def edge(vertices, from_id, to_id):
    f = vertices.setdefault(from_id, Vertex(from_id))
    t = vertices.setdefault(to_id, Vertex(to_id))
    f.add_edge(t)

def pagerank_test(vertices, num_vertices): # [Vertex] -> 1D Matrix
    """Computes the pagerank vector associated to vertices, using a
    standard matrix-theoretic approach to computing pagerank.  This is
    used as a basis for comparison."""
    I = mat(eye(num_vertices))
    G = zeros((num_vertices,num_vertices))
    for vertex in vertices:
        num_out_vertices = len(vertex.out_vertices)
        for out_vertex in vertex.out_vertices:
            G[out_vertex.id,vertex.id] = 1.0/num_out_vertices
    P = (1.0/num_vertices)*mat(ones((num_vertices,1)))
    return 0.15*((I-0.85*G).I)*P

def print_prs(prs): # 1DMatrix -> out
    v_prs = zip(range(0, prs.size), prs.flat) # [vertex_id, page_rank]
    for pr in v_prs:
        print(pr[0] + 1, pr[1])

def set_init_value(vertices, num_vertices): # [Vertex] ->
    """Initializes the vertex values to 1/num_vertices"""
    for v in vertices:
        v.value = 1.0/num_vertices

if __name__ == "__main__":
    main()

