use strict;
my @c = (
    "This program is for humans",
    "Still in Mama's arms",
    "Preschool Maniac",
    "Elementary school",
    "Middle school",
    "High school",
    "College",
    "Working for the man",
    "The Golden Years"
);
my @a = ( 0, 3, 5, 12, 15, 19, 23, 66, 101 );
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    my $i = 0;
    $i++ while $i < 9 && $line >= @a[$i];
    printf( "%s\n", @c[ $i % scalar @c ] );
}
close(INFILE);
