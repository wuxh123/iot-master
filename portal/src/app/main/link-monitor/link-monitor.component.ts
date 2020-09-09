import {Component, ElementRef, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {ApiService} from '../../api.service';

@Component({
  selector: 'app-link-monitor',
  templateUrl: './link-monitor.component.html',
  styleUrls: ['./link-monitor.component.scss']
})
export class LinkMonitorComponent implements OnInit, OnDestroy {
  @ViewChild('contentRecv')
  contentRecv: ElementRef;

  @ViewChild('contentSend')
  contentSend: ElementRef;

  id: number;
  link: any;
  ws: WebSocket;
  interval: any;

  text = '';
  dataRecv = [];
  dataSend = [];

  cacheSizeRecv = 500;
  cacheSizeSend = 500;

  constructor(private routeInfo: ActivatedRoute, private as: ApiService) {
    this.id = this.routeInfo.snapshot.params.id;
    this.load();
  }

  ngOnInit(): void {
  }

  ngOnDestroy(): void {
    this.ws.close(1000, 'exit');
    clearInterval(this.interval);
  }

  startHearbeat(): void {
    this.interval = setInterval(() => {
      this.ws.send(JSON.stringify({type: 'ping'}));
    }, 10000);
  }

  load(): void {
    this.as.get('link/' + this.id).subscribe(res => {
      this.link = res.data;
      // TODO 检查在线状态
      this.monitor();
    });
  }

  loadStatus(): void {
    // 在线，monitor
  }

  send(): void {
    console.log('send', this.text);
    this.ws.send(JSON.stringify({type: 'send', data: this.text}));
  }

  monitor(): void {
    this.ws = new WebSocket('ws://127.0.0.1:8080/api/monitor/' + this.link.channel_id + '/' + this.id);

    this.ws.onopen = e => {
      console.log('Connection open ...');
      // ws.send("{}");
      this.startHearbeat();
    };

    this.ws.onmessage = e => {
      console.log('Recv: ' + e.data);
      const obj = JSON.parse(e.data);
      switch (obj.type) {
        case 'recv':
          this.dataRecv.push(obj);
          if (this.dataRecv.length > this.cacheSizeRecv) {
            this.dataRecv.splice(0, this.dataRecv.length - this.cacheSizeRecv);
          }
          this.contentRecv.nativeElement.scrollTo(0, this.contentRecv.nativeElement.scrollHeight);
          break;
        case 'send':
          this.dataSend.push(obj);
          if (this.dataSend.length > this.cacheSizeSend) {
            this.dataSend.splice(0, this.dataSend.length - this.cacheSizeSend);
          }
          this.contentSend.nativeElement.scrollTo(0, this.contentSend.nativeElement.scrollHeight);
          break;
      }
    };

    this.ws.onclose = e => {
      console.log('Connection closed.');
    };
  }

}
