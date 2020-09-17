import {Injectable} from '@angular/core';

import * as mqtt from 'mqtt';
import {Observable, Subject, Subscription, Unsubscribable, using, merge} from "rxjs";
import {filter, publish, refCount} from "rxjs/operators";


@Injectable({
  providedIn: 'root'
})
export class MqttService {
  client: mqtt.MqttClient;
  messages: Subject<any> = new Subject();

  topics: { [key: string]: Observable<any> } = {};

  constructor() {
    const client = this.client = mqtt.connect('ws://127.0.0.1:8080/api/mqtt');
    client.on('connect', data => {
      console.log('mqtt connect', data);
    });
    client.on('message', (topic, message, packet) => {
      console.log('mqtt message', topic, message);
      if (packet.cmd === 'publish') {
        this.messages.next(packet);
      }
    });
    client.on('close', () => {
      console.log('mqtt close');
    });
    client.on('offline', () => {
      console.log('mqtt offline');
    });
    client.on('disconnect', (data) => {
      console.log('mqtt disconnect');
    });
    client.on('error', (err) => {
      console.log('mqtt error', err);
    });
  }

  match(sub: string, topic: string): boolean {
    const fs = sub.split('/');
    const ts = topic.split('/');
    let i = 0;
    for (; i < fs.length && i < ts.length; i++) {
      if (fs[i] === '#') {
        return true;
      } else if (fs[i] === '+' || fs[i] === ts[i]) {
        // continue;
      } else {
        return false;
      }
    }
    return i === fs.length && i === ts.length;
  }


  subscribe(filterString: string): Observable<any> {
    if (!this.topics[filterString]) {
      const rejected: Subject<any> = new Subject();
      this.topics[filterString] = using(
        // topics: Do the actual ref-counting MQTT subscription.
        // refcount is decreased on unsubscribe.
        () => {
          const subscription: Subscription = new Subscription();
          this.client.subscribe(filterString, {qos: 1}, (err) => {
          });
          subscription.add(() => {
            delete this.topics[filterString];
            this.client.unsubscribe(filterString);
          });
          return subscription;
        },
        // observableFactory: Create the observable that is consumed from.
        // This part is not executed until the Observable returned by
        // `observe` gets actually subscribed.
        (subscription: Unsubscribable | void) => merge(rejected, this.messages))
        .pipe(
          filter((msg: any) => this.match(filterString, msg.topic)),
          publish(),
          refCount()
        ) as Observable<Buffer>;
    }
    return this.topics[filterString];
  }

}
