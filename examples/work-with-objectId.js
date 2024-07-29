import xk6_mongodb from 'k6/x/mongodb'

const connection = xk6_mongodb.connect('mongodb://localhost:27017')

const doc = {
  name: 'k6',
  ref_id: xk6_mongodb.newId(),
}

export default function () {
  connection.insert('k6', 'k6', doc)
}
