import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {TabRef} from '../tabs/tabs.component';
import {ApiService} from '../../api.service';

@Component({
  selector: 'app-element-detail',
  templateUrl: './element-detail.component.html',
  styleUrls: ['./element-detail.component.scss']
})
export class ElementDetailComponent implements OnInit {

  element: any = {};
  id = 0;

  constructor(private as: ApiService, private routeInfo: ActivatedRoute, private tab: TabRef) {
    this.id = routeInfo.snapshot.params.id;
    tab.name = '元件详情';
  }

  ngOnInit(): void {
    this.as.get('element/' + this.id).subscribe(res => {
      if (res.ok) {
        this.element = res.data;
        this.tab.name = '元件【' + this.element.name + '】';
      }
      // console.log(res);
    });
  }

}
