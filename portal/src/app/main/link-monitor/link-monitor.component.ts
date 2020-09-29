import {Component, ElementRef, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {ApiService} from '../../api.service';
import {MqttService} from '../../mqtt.service';

@Component({
  selector: 'app-link-monitor',
  templateUrl: './link-monitor.component.html',
  styleUrls: ['./link-monitor.component.scss']
})
export class LinkMonitorComponent implements OnInit, OnDestroy {
  title = '连接监控';


  @ViewChild('contentRecv')
  contentRecv: ElementRef;

  @ViewChild('contentSend')
  contentSend: ElementRef;

  id: number;
  link: any;

  isHex = false;

  text = '';
  dataRecv = [];
  dataSend = [];

  cacheSizeRecv = 500;
  cacheSizeSend = 500;

  recvSub: any;
  sendSub: any;

  constructor(private routeInfo: ActivatedRoute, private as: ApiService, private mqtt: MqttService) {
    this.id = this.routeInfo.snapshot.params.id;
    this.load();
  }

  ngOnInit(): void {

  }

  ngOnDestroy(): void {
    this.recvSub.unsubscribe();
    this.sendSub.unsubscribe();
  }

  hex_to_buffer(hex: string): Buffer {
    hex = hex.replace(/\s*/g, '');
    const arr = [];
    for (let i = 0; i < hex.length; i += 2) {
      arr.push(hex.substr(i, 2));
    }
    const hexes = arr.map(el => parseInt(el, 16));
    return Buffer.from(new Uint8Array(hexes));
  }

  buffer_to_hex(buf): string {
    const arr = Array.prototype.slice.call(buf);
    return arr.map(el => Number(el).toString(16)).join(' ');
  }

  subscribe(): void {
    this.recvSub = this.mqtt.subscribe('/' + this.link.channel_id + '/' + this.id + '/recv').subscribe(packet => {
      this.dataRecv.push({
        data: this.buffer_to_hex(packet.payload),
        time: new Date(),
      });
      if (this.dataRecv.length > this.cacheSizeRecv) {
        this.dataRecv.splice(0, 1);
      }
      this.contentRecv.nativeElement.scrollTo(0, this.contentRecv.nativeElement.scrollHeight);
    });

    this.sendSub = this.mqtt.subscribe('/' + this.link.channel_id + '/' + this.id + '/send').subscribe(packet => {
      this.dataSend.push({
        data: this.buffer_to_hex(packet.payload),
        time: new Date(),
      });
      if (this.dataSend.length > this.cacheSizeSend) {
        this.dataSend.splice(0, 1);
      }
      this.contentSend.nativeElement.scrollTo(0, this.contentSend.nativeElement.scrollHeight);
    });
  }

  load(): void {
    this.as.get('link/' + this.id).subscribe(res => {
      this.link = res.data;

      this.subscribe();
    });
  }

  loadStatus(): void {
    // 在线，monitor
  }

  send(): void {
    //console.log('send', this.text);
    let content: any = this.text;
    // 转换十六进制
    if (this.isHex) {
      content = this.hex_to_buffer(this.text);
    }
    this.mqtt.publish('/' + this.link.channel_id + '/' + this.id + '/transfer', content);
  }

}
