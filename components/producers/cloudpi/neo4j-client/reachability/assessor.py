from neo4j import GraphDatabase

uri = "bolt://localhost:7687"

# Create a Neo4j driver instance
driver = GraphDatabase.driver(uri)

def get_snapshosts_exposed():
    query = """
    MATCH (s:EBSSnapshot)-[:CREATED_FROM]->(v:EBSVolume),
    (v)-[:ATTACHED_TO]->(i:EC2Instance{exposed_internet:true})
    RETURN i.id AS instance_id, v.id AS volume_id, s.id AS snapshot_id, s.lastupdated AS last_updated  LIMIT 25
    """
    with driver.session() as session:
        result = session.run(query)
        snapshots = [
            {
                "instance_id": record["instance_id"], 
                "volume_id": record["volume_id"],
                "snapshot_id": record["snapshot_id"], 
                "last_updated": record["last_updated"],
            } 
        for record in result]
        return snapshots
