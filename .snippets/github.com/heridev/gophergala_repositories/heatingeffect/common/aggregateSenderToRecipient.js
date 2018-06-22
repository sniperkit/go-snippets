db.notices.aggregate([
{
	$match: {
		sender_name: { $exists: true},
		recipient_name: {$exists: true}
	}
},
{
    $group: {
        _id: {
		sender: "$sender_name",
		recipient: "$recipient_name"
	},
	notices: { "$sum": 1 },
    }
},
{
	$out: "notices_sendto_stat"
}
])