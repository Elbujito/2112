use redis:: { Commands, RedisResult };
use serde_json:: Value;

pub fn publish_to_redis(channel: & str, message: & Value, redis_client: & redis:: Client) -> RedisResult < () > {
    let mut con = redis_client.get_connection() ?;
    let _: () = con.publish(channel, message.to_string()) ?;
    Ok(())
}

pub fn store_to_redis(key: & str, data: & Value, score: f64, redis_client: & redis:: Client) -> RedisResult < () > {
    let mut con = redis_client.get_connection() ?;
    let _: () = con.zadd(key, data.to_string(), score) ?;
    Ok(())
}
