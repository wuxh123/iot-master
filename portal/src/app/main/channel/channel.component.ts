import {Component, OnInit} from '@angular/core';
import {ApiService} from "../../api.service";

@Component({
  selector: 'app-channel',
  templateUrl: './channel.component.html',
  styleUrls: ['./channel.component.scss']
})
export class ChannelComponent implements OnInit {

  channels: [];

  constructor(private as: ApiService) {
  }

  ngOnInit(): void {
    this.load();
  }

  load(): void {
    this.as.get('channels').subscribe(res => {
      if (res.ok) {
        this.channels = res.data;
      }
    });
  }

  create(): void {

  }

}
