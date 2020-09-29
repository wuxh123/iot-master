import {Component, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';

@Component({
  selector: 'app-dash',
  templateUrl: './dash.component.html',
  styleUrls: ['./dash.component.scss']
})
export class DashComponent implements OnInit {
  title = '仪表盘';

  constructor(private as: ApiService) {
  }

  ngOnInit(): void {

  }

}
