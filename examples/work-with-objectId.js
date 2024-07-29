import xk6_mongodb from 'k6/x/mongodb'

const connection = xk6_mongodb.connect('mongodb://localhost:27017')

let doc = {
  name: 'k6',
  ref_id: {
    $oid: '66a6f37b369b0796f6f7f840'
  }
}
doc = connection.transformDoc(doc)

export default function () {
  connection.insert('k6', 'k6', doc)
}
