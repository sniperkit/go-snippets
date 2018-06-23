"""This file runs the prototype"""

import data
import ranking
import composition_functions


def valid_features(query):
    valid = True
    if len(query) % 2 != 0:
        valid = False
    for key in dico:
        if key not in data_source.features:
            valid = False
    return valid


if __name__ == "__main__":
    import sys
    from pprint import pprint

    seed = None
    if len(sys.argv) > 1:
        seed = sys.argv[1]

    data_source = data.RandomData(seed)
    ranking = ranking.Ranking(data_source, composition_functions.weighted_sum)

    how_to = ('\nTo write your query, you can mention the feature(s) you are '
              'interested\nin and give them some weight, or importance.\n'
              'It will return the most skilled users combining those skills, '
              'and their scores,\naccording to DevMine\'s crazy algorithm.\n\n'
              'For example, the query: python 5 superpower 10\n'
              'will give you:\n\n')
    example = """
        [('Kevin', 11.240153556173475),
         ('Clément', 8.545413290545826),
         ('Laurent', 7.661851902424877),
         ('Frederik', 7.547396298747576),
         ('Robin', 7.060494410357596),
         ('Daniel', 5.780794009967852),
         ('Xuan', 4.955532838858104)]\n"""
    print(how_to, example)

    print('\nHere are the features you can query over:\n')
    for features in data_source.features:
        print(features)

    print('\nEnter "quit" to stop your search of amazing developers.\n')

    query = ''
    query = input('Enter your query: \n\n')

    while query != 'quit':

        query = query.split()
        dico = dict()

        try:
            for i in range(len(query)):
                if i % 2 == 0:
                    dico[query[i].lower()] = int(query[i+1])

            if valid_features(query):
                pprint(ranking.rank_all(dico)[:10])
            else:
                print('Invalid query')
        except (TypeError, IndexError):
            print('Invalid query')

        query = input('\n\nEnter your query: ')
