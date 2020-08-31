import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-channel-detail',
  templateUrl: './channel-detail.component.html',
  styleUrls: ['./channel-detail.component.scss']
})
export class ChannelDetailComponent implements OnInit {

  @Input() channel = {};

  constructor() { }

  ngOnInit(): void {
  }

}
