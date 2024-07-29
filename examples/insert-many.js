import xk6_mongodb from 'k6/x/mongodb'

const connection = xk6_mongodb.connect('mongodb://localhost:27017')

export default function () {
  connection.insertMany('k6', 'k6', [{ name: 'k6', n: 1 }, { name: 'k6', n: 2 }])
}
